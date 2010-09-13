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

package python

import "bytes"
import "encoding/binary"

const (
    NOP = iota          // 0 - 15 are "special" instructions
    NEW        
    LEN
)

const (    
    LOAD = 16 + iota    // 16-32 are immediate-mode instructions
    BIND
    BOXI    
    BOXL
    BOXF
    BOXS
    BOXB
    UNBOXI
    UNBOXL
    UNBOXF
    UNBOXS
    UNBOXB
)

const ( 
    INDEX = 33 + iota   // 33-63 are register 3-code instructions op (src1, src2, dst)
    SPILL
    FILL
    SET
    GET
    ADD
    SUB
    MUL
    DIV
    MOD
)

// A code stream contains all the code for one module
type CodeStream struct {
    *bytes.Buffer
        
    Strings         map[string]uint16
    StringCounter   uint16
    
    Locals          map[uint16]*Object
    Globals         map[uint16]*Object        
}

func (s *CodeStream) WriteLoad(name string, register uint32) {
    var instruction uint32
    var value uint16
    var present bool
    
    value, present = s.Strings[name]
    
    if !present {        
        value = s.StringCounter
        s.StringCounter++
    }

    instruction = LOAD | (uint32(value) << immediate_val_shift) | (register << imm_target_reg_shift)
    binary.Write(s, binary.LittleEndian, instruction)    
}

func (s *CodeStream) WriteBind(name string, register uint32) {
    var instruction uint32 = 0
    var value uint16
    var present bool
    
    value, present = s.Strings[name]
    
    if !present {        
        value = s.StringCounter
        s.StringCounter++
    }

    instruction = BIND | (uint32(value) << immediate_val_shift) | (register << imm_target_reg_shift)
    binary.Write(s, binary.LittleEndian, instruction)    
}
