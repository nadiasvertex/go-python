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

  The VM package provides a Go implementation of the Python virtual machine.  The VM
  is implemented as a register-based system with tracing and JIT hooks.   
  
*/

package vm

import "bytes"

const (
    NOP,
    DEREF,
    BIND,
    NEW,
    INDEX,
    SET,
    GET,
    BOX,
    UNBOX,
    ADD,
    SUB,
    MUL,
    DIV,
    MOD,
    LEN    
)

// A code stream contains all the code for one module.
type CodeStream struct {
    *bytes.Buffer
        
    Strings         map[string]int
    StringCounter   int
    
    Locals          map[int]*Object
    Globals         map[int]*Object        
}

func (s *CodeStream) WriteDeref(name string, register byte) {
    var value int
    var present bool
    
    if value, present := s.Strings[name]; !present {        
        value = s.StringCounter
        s.StringCounter++
    }

    s.Write(DEREF)
    s.Write(value)
    s.Write(register)
}