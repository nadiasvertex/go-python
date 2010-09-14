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

import (
        "testing"            
        //"encoding/binary"
)

var sample_instructions = []int32{0x0003a010, 0x05060708}

func TestEncodeInstructions(t *testing.T) {
    
    s := new (CodeStream)
    s.Init()

    s.WriteLoad("a", 3)
    s.WriteBind("b", 5)    
  
}
