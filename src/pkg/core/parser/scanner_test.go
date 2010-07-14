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
    token{Integer, "0b11"},
    token{Integer, "010"},
    token{Integer, "19"},
    token{Integer, "0xAc"},    
}

func makeSource(pattern string) *bytes.Buffer {
    var buf bytes.Buffer
    for _, k := range tokenList {
        fmt.Fprintf(&buf, pattern, k.text)
    }
    return &buf
}

func TestScanNumber(t *testing.T) {
    s := new(Scanner).Init(makeSource("%s\n"))
    
    tok := s.Scan()        
    
    if (tok!=Integer) {
        t.Errorf("Expected an integer but got something else.")
    }   
}

