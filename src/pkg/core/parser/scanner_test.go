/* Copyright 2010 Christopher Nelson

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
*/

package parser

import ( 
    "bytes";
    "fmt";
    "testing"
)

type token struct {
    tok  int
    text string
}

var tokenList = []token{
    token{Integer, "0b10"},
    token{Integer, "01234567"},
    token{Integer, "1234567890"},
    token{Integer, "0xabcdef0123456789FEDCBA"},  
    
    token{Indent, "  "},
    token{Dedent, " "},
    
    token{Identifier, "print"},
    token{Identifier, "call_forward"},
    token{Identifier, "Parser5"},
    
    token{String, "\"test\""},    
    token{String, "'test2'"}, 
    token{String, "\"\"\"test\nand\ntest\"\"\""},    
    token{String, "'''test2\nand\ntest2'''"},
    token{String, "r\"raw_test\""},     
    token{String, "r'raw_test2'"},
}

func makeSource(pattern string) *bytes.Buffer {
    var buf bytes.Buffer
    for _, k := range tokenList {
        fmt.Fprintf(&buf, pattern, k.text)
    }
    return &buf
}

func TestScanTokens(t *testing.T) {
    s := new(Scanner).Init(makeSource("%s\n"))
    
    tok := s.Scan()        
    
    for _, k := range tokenList {
        // Ignore EOL that happens after each scan.
        if tok == EOL {
            tok=s.Scan()
        }
                       
        if tok != k.tok {
            t.Fatalf("%d:%d Expected token type '%s' but got '%s' for '%s' (token text='%s')", s.line, s.column, tokenString[k.tok], tokenString[tok], k.text, s.TokenText())
        } else if k.text != s.TokenText() {
            t.Errorf("%d:%d Expected '%s' but got '%s' for token '%s'", s.line, s.column, k.text, s.TokenText(), tokenString[tok])
        }        
    
        tok = s.Scan()    
    }
       
}

