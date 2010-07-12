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
*/

package parser

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
    Indent
    Dedent
    Identifier
    Integer
    Long
    Float    
    Imaginary
    String
)

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

    // The Mode field controls which tokens are recognized. For instance,
    // to recognize Ints, set the (1<<-Int) bit in Mode. The field may be
    // changed at any time.
    Mode uint

    // The Whitespace field controls which characters are recognized
    // as white space. To recognize a character ch <= ' ' as white space,
    // set the ch'th bit in Whitespace (the Scanner's behavior is undefined
    // for values ch > ' '). The field may be changed at any time.
    Whitespace uint64

    // Current token position. The Offset, Line, and Column fields
    // are set by Scan(); the Filename field is left untouched by the
    // Scanner.
    Position
}

// Init initializes a Scanner with a new source and returns itself.
// Error is set to nil, ErrorCount is set to 0, Mode is set to GoTokens,
// and Whitespace is set to GoWhitespace.
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

    // initialize token text buffer
    s.tokPos = -1

    // initialize one character look-ahead
    s.ch = s.next()

    // initialize public fields
    s.Error = nil
    s.ErrorCount = 0
    s.Mode = GoTokens
    s.Whitespace = GoWhitespace

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


func (s *Scanner) scanIdentifier() int {
    ch := s.next() // read character after first '_' or letter
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
			case 'o', 'O':
				ch = s.next()
				for isOctDigit(ch) {
					ch = s.next()
				}
				return Integer, ch
			
			case 'x', 'X':
				ch = s.next()
				for isHexDigit(ch) {
					ch = s.next()
				}
				return Integer, ch
			
			case 'b', 'B':
				ch = s.next()
				for isBinDigit(ch) {
					ch = s.next()
				}				
				return Integer, ch
		}
	
	}
}
