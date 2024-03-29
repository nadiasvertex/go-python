/* 
   Copyright 2010 Christopher Nelson

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
   --------------------------------------------------------------------

   The parser package implements a simple library for parsing EBNF
   grammars.
   
   The scanner, lexer, and parser are all implemented together for
   efficiency.  Much of the scanner was happily stolen from the Go scanner package
   and reworked to be specific to Python.
   
   In particular, this scanner completely re-implements the Scan() and scanXXXX() functions,
   removes mode and configurable whitespace, tokenizes the newline character, and maintains
   an indent/dedent stack - which also appear as tokens.
*/

package python

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "unicode"
    "utf8"
)

// A source position is represented by a Position value.
// A position is valid if Line > 0.
type Position struct {
    Filename string // filename, if any
    Offset   int    // byte offset, starting at 0
    Line     int    // line number, starting at 1
    Column   int    // column number, starting at 0 (character count per line)
}

// IsValid returns true if the position is valid.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

func (pos Position) String() string {
    s := pos.Filename
    if pos.IsValid() {
        if s != "" {
            s += ":"
        }
        s += fmt.Sprintf("%d:%d", pos.Line, pos.Column)
    }
    if s == "" {
        s = "???"
    }
    return s
}

const (
    EOF = -(iota + 1)
    EOL 
    Indent
    Dedent
    Identifier
    Integer
    Long
    Float    
    Imaginary
    String
    Comment
)

var tokenString = map[int]string{
    EOF:        "EOF",
    EOL:        "EOL",
    Indent:     "Indent",
    Dedent:     "Dedent",
    Identifier: "Identifier",
    Integer:    "Integer",
    Float:      "Float",
    Long:       "Long",
    String:     "String",
    Imaginary:  "Imaginary",
    Comment:    "Comment",
}

const bufLen = 1024 // at least utf8.UTFMax

// A Scanner implements reading of Unicode characters and tokens from an io.Reader.
type Scanner struct {
    // Input
    src io.Reader

    // Source buffer
    srcBuf [bufLen + 1]byte // +1 for sentinel for common case of s.next()
    srcPos int              // reading position (srcBuf index)
    srcEnd int              // source end (srcBuf index)

    // Source position
    srcBufOffset int // byte offset of srcBuf[0] in source
    line         int // newline count + 1
    column       int // character count on line
    
    // Some state necessary for Python-esque token scanning
    isNewline    bool     // if we just returned an EOL token, this is true.
    indentStack [1024]int // the indent stack, keeps track of the various indent levels
    indentPos   int       // the stack pointer for the indent. indicates top of stack.

    // Token text buffer
    // Typically, token text is stored completely in srcBuf, but in general
    // the token text's head may be buffered in tokBuf while the token text's
    // tail is stored in srcBuf.
    tokBuf bytes.Buffer // token text head that is not in srcBuf anymore
    tokPos int          // token text tail position (srcBuf index)
    tokEnd int          // token text tail end (srcBuf index)

    // One character look-ahead
    ch int // character before current srcPos

    // Error is called for each error encountered. If no Error
    // function is set, the error is reported to os.Stderr.
    Error func(s *Scanner, msg string)

    // ErrorCount is incremented by one for each error encountered.
    ErrorCount int
        
    // Current token position. The Offset, Line, and Column fields
    // are set by Scan(); the Filename field is left untouched by the
    // Scanner.
    Position
}

// Init initializes a Scanner with a new source and returns itself.
// Error is set to nil, and ErrorCount is set to 0.
func (s *Scanner) Init(src io.Reader) *Scanner {
    s.src = src

    // initialize source buffer
    s.srcBuf[0] = utf8.RuneSelf // sentinel
    s.srcPos = 0
    s.srcEnd = 0

    // initialize source position
    s.srcBufOffset = 0
    s.line = 1
    s.column = 0
    
    // initialize indent tracker
    s.isNewline = true
    s.indentPos = 0

    // initialize token text buffer
    s.tokPos = -1

    // initialize one character look-ahead
    s.ch = s.next()

    // initialize public fields
    s.Error = nil
    s.ErrorCount = 0
    
    return s
}


// next reads and returns the next Unicode character. It is designed such
// that only a minimal amount of work needs to be done in the common ASCII
// case (one test to check for both ASCII and end-of-buffer, and one test
// to check for newlines).
func (s *Scanner) next() int {
    ch := int(s.srcBuf[s.srcPos])

    if ch >= utf8.RuneSelf {
        // uncommon case: not ASCII or not enough bytes
        for s.srcPos+utf8.UTFMax > s.srcEnd && !utf8.FullRune(s.srcBuf[s.srcPos:s.srcEnd]) {
            // not enough bytes: read some more, but first
            // save away token text if any
            if s.tokPos >= 0 {
                s.tokBuf.Write(s.srcBuf[s.tokPos:s.srcPos])
                s.tokPos = 0
            }
            // move unread bytes to beginning of buffer
            copy(s.srcBuf[0:], s.srcBuf[s.srcPos:s.srcEnd])
            s.srcBufOffset += s.srcPos
            // read more bytes
            i := s.srcEnd - s.srcPos
            n, err := s.src.Read(s.srcBuf[i:bufLen])
            s.srcEnd = i + n
            s.srcPos = 0
            s.srcBuf[s.srcEnd] = utf8.RuneSelf // sentinel
            if err != nil {
                if s.srcEnd == 0 {
                    return EOF
                }
                s.error(err.String())
                break
            }
        }
        // at least one byte
        ch = int(s.srcBuf[s.srcPos])
        if ch >= utf8.RuneSelf {
            // uncommon case: not ASCII
            var width int
            ch, width = utf8.DecodeRune(s.srcBuf[s.srcPos:s.srcEnd])
            if ch == utf8.RuneError && width == 1 {
                s.error("illegal UTF-8 encoding")
            }
            s.srcPos += width - 1
        }
    }

    s.srcPos++
    s.column++
    switch ch {
    case 0:
        // implementation restriction for compatibility with other tools
        s.error("illegal character NUL")
    case '\n':
        s.line++
        s.column = 0
    }

    return ch
}


// Next reads and returns the next Unicode character.
// It returns EOF at the end of the source. It reports
// a read error by calling s.Error, if set, or else
// prints an error message to os.Stderr. Next does not
// update the Scanner's Position field; use Pos() to
// get the current position.
func (s *Scanner) Next() int {
    s.tokPos = -1 // don't collect token text
    ch := s.ch
    s.ch = s.next()
    return ch
}


// Peek returns the next Unicode character in the source without advancing
// the scanner. It returns EOF if the scanner's position is at the last
// character of the source.
func (s *Scanner) Peek() int {
    return s.ch
}

func (s *Scanner) error(msg string) {
    s.ErrorCount++
    if s.Error != nil {
        s.Error(s, msg)
        return
    }
    fmt.Fprintf(os.Stderr, "%s: %s", s.Position, msg)
}

func (s *Scanner) scanIdentifier(ch int) int {    
    for ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) {
        ch = s.next()
    }
    return ch
}

func isBinDigit(ch int) bool {
	switch ch {
		case '0', '1':
			return true
	}	
	return false
}

func isOctDigit(ch int) bool {
	switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7':
			return true
	}	
	return false
}

func isDecDigit(ch int) bool {
    switch ch {
        case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
            return true
    }   
    return false
}

func isHexDigit(ch int) bool {
	switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'a', 'B', 'b', 'C', 'c', 'D', 'd', 'E', 'e', 'F', 'f':
			return true
	}	
	return false
}

func (s *Scanner) scanNumber(ch int) (int, int) {
	// Not a decimal number
	if ch == '0' {
		ch = s.next()
		switch ch {
		    
		    // Scan hex int
			case 'x', 'X':
				ch = s.next()
				for isHexDigit(ch) {
					ch = s.next()
				}				
			
			// Scan binary int
			case 'b', 'B':
				ch = s.next()
				for isBinDigit(ch) {
					ch = s.next()
				}
			
			// Scan dec int	
		    default:
		        ch = s.next()
                for isOctDigit(ch) {
                    ch = s.next()
                }               
            
		}	
	} else {
        // Decimal number	
        for isDecDigit(ch) {
            ch = s.next()
        }
    }
	
	return Integer, ch	
}

func (s *Scanner) scanString(quote int) (n int) {
    multiline := false
    ch := s.next() // read character after quote
    
    // Handle multiline strings
    if ch == quote && s.Peek() == quote {
        multiline = true
        ch = s.next()
        ch = s.next()
    }
    for ch != quote {
        if (!multiline && ch == '\n') || ch < 0 {
            s.error("string literal not terminated\n")
            return
        }
        if ch == '\\' {
            ch = s.next() //s.scanEscape(quote)
        } else {
            ch = s.next()
        }
        n++
    }

    // Consume the extra quote characters when scanning
    // multiline Python strings.    
    if multiline {
        ch = s.next()
        ch = s.next()
    }
    
    return
}


// Scan reads the next token or Unicode character from source and returns it.
// It returns EOF at the end of the source. It reports scanner errors (read and
// token errors) by calling s.Error, if set; otherwise it prints an error message
// to os.Stderr.
func (s *Scanner) Scan() int {
    ch := s.ch

    // reset token text position
    s.tokPos = -1

redo:
    // skip white space
    if !s.isNewline {
        for ch == ' ' || ch == '\t' {
            ch = s.next()
        }
    }
    
    // start collecting token text
    s.tokBuf.Reset()
    s.tokPos = s.srcPos - 1

    // set token position
    s.Offset = s.srcBufOffset + s.tokPos
    s.Line = s.line
    s.Column = s.column

    // determine token value
    tok := ch
    switch {
        case unicode.IsLetter(ch) || ch == '_':            
            scan_identifier := true
            
            // Handle raw strings, which look like identifiers at the beginning.
            if (ch == 'r' || ch=='u') {
                ch = s.next()
                if ch == '"' || ch == '\'' {
                    scan_identifier = false
                    s.scanString(ch)
                    tok = String
                    ch = s.next()
                }
            } 
            
            // Handle identifiers
            if scan_identifier {                 
                tok = Identifier
                ch = s.scanIdentifier(ch)
            }
          
        case isDecDigit(ch):        
            tok, ch = s.scanNumber(ch)
            
        case ch == '\\':
            // Handle explicit line joining.            
            ch = s.next()
            for ch =='\r' || ch == '\n' {
                ch = s.next()
            }
            
            goto redo
                
        case ch == '\r' || ch == '\n':
            // Handle end of line reporting
            tok = EOL
            // Check for /r/n or just /r line endings
            if ch=='\r' {
                ch = s.next()
                if ch=='\n' {
                    ch = s.next()
                }
            }       
            
            ch = s.next()
            
        case ch == ' ' || ch == '\t':
            // handle indent / dedent    
            indent_length := 0
            for ch == ' ' || ch == '\t' {
                switch ch {
                    case  ' ': indent_length += 1                       // increase indent by 1
                    case '\t': indent_length = ((indent_length/8)+1)*8  // pad indent to nearest multiple of 8 (Python lex spec rule.)
                }
                
                ch = s.next()
            }
            
            // Figure out if we should emit an indent, dedent, or
            // nothing.  If the indentation level hasn't changed
            // we ignore the whitespace.
            switch {
                case indent_length > s.indentStack[s.indentPos]: 
                    tok = Indent
                    s.indentPos++
                    s.indentStack[s.indentPos] = indent_length
                    
                case indent_length < s.indentStack[s.indentPos]: 
                    tok = Dedent
                    s.indentPos++
                    s.indentStack[s.indentPos] = indent_length                
                    
                default:
                    goto redo            
            }             
                        
            
        default:
            switch ch {      
                case '"', '\'':
                    s.scanString(ch)
                    tok = String
                    ch = s.next()
                default:
                    ch = s.next()
            }
    }

    // end of token textindent_length += 1
    s.tokEnd = s.srcPos - 1

    // process newline
    s.isNewline = (tok == EOL)    

    s.ch = ch
    return tok
}

// Position returns the current source position. If called before Next()
// or Scan(), it returns the position of the next Unicode character or token
// returned by these functions. If called afterwards, it returns the position
// immediately after the last character of the most recent token or character
// scanned.
func (s *Scanner) Pos() Position {
    return Position{
        s.Filename,
        s.srcBufOffset + s.srcPos - 1,
        s.line,
        s.column,
    }
}


// TokenText returns the string corresponding to the most recently scanned token.
// Valid after calling Scan().
func (s *Scanner) TokenText() string {
    if s.tokPos < 0 {
        // no token text
        return ""
    }

    if s.tokEnd < 0 {
        // if EOF was reached, s.tokEnd is set to -1 (s.srcPos == 0)
        s.tokEnd = s.tokPos
    }

    if s.tokBuf.Len() == 0 {
        // common case: the entire token text is still in srcBuf
        return string(s.srcBuf[s.tokPos:s.tokEnd])
    }

    // part of the token text was saved in tokBuf: save the rest in
    // tokBuf as well and return its content
    s.tokBuf.Write(s.srcBuf[s.tokPos:s.tokEnd])
    s.tokPos = s.tokEnd // ensure idempotency of TokenText() call
    return s.tokBuf.String()
}
