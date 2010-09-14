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
)

func checkIntResult(t *testing.T, m *Machine, register int, wanted Object, message string) {
    if test_value, ok := m.Register[register].(*IntObject); !ok {
        t.Errorf("failure dispatching '%v' (register %v has incorrect type: '%v')\n", message, register, m.Register[register])        
    } else {
        if m.Register[register] != wanted {
            t.Errorf("failure dispatching '%v', (register %v has incorrect value '%v')\n", message, register, test_value.AsInt())            
        }
    }
    
    return
}

func checkIntValueResult(t *testing.T, m *Machine, register int, wanted int, message string) {
    if test_value, ok := m.Register[register].(*IntObject); !ok {
        t.Errorf("failure dispatching '%v' (register %v has incorrect type: '%v')\n", message, register, m.Register[register])        
    } else {
        if m.Register[register].AsInt() != wanted {
            t.Errorf("failure dispatching '%v', (register %v has incorrect value '%v' wanted '%v')\n", message, register, test_value.AsInt(), wanted)            
        }
    }
    
    return
}


func TestDispatchInstructions(t *testing.T) {
    
    s := new (CodeStream)
    s.Init()
    
    m := new (Machine)
    
    io1 := new(IntObject)
    io1.Value = 10
            
    s.BindLocal("a", io1)    

    s.WriteLoad("a", 3, false, 0)
    s.WriteBind("b", 3, false, 0)
    s.WriteLoad("b", 4, false, 0)
    s.WriteAdd(3,4,5,false,0)
    
    // Test the Load
    m.Dispatch(s)
    checkIntResult(t, m, 3, io1, "LOAD 1, r3")
    
    // Test bind and reload
    m.Dispatch(s)
    m.Dispatch(s)
    checkIntResult(t, m, 4, io1, "LOAD 2, r4")
    
    // Test add
    m.Dispatch(s)    
    checkIntValueResult(t, m, 5, 20, "ADD r3, r4, r5")           
    
}
