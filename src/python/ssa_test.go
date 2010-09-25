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

  Test for the static single assignment subsystem.   
  
*/

package python

import (   
        "big"     
        "testing"            
)


func TestWriteElements(t *testing.T) {    
    ctx := new (SsaContext)
    ctx.Init()
    
    el := new (SsaElement)
    
    for i:=0; i<256; i++ {
        if ctx.Write(el) != i {
            t.Errorf("Write returned an incorrect value as the new id.\n")  
        }
    }    
}

func TestLoadInt(t *testing.T) {    
    ctx := new (SsaContext)
    ctx.Init()
        
    for i:=0; i<256; i++ {
        if ctx.LoadInt(big.NewInt(int64(i))) != i {
            t.Errorf("NameInt returned an incorrect value as the new id.\n")  
        }
    }
    
    new_int := big.NewInt(1000)
    new_int_idx := ctx.LoadInt(new_int)
    
    for i:=0; i<256; i++ {
        if ctx.LoadInt(new_int) != new_int_idx {
            t.Errorf("NameInt returned an incorrect value as the new id for an identical name.\n")  
        }
    }    
}

func TestEval(t *testing.T) {    
    ctx := new (SsaContext)
    ctx.Init()
    
    left_int  := big.NewInt(1000)
    right_int := big.NewInt(500)
       
    // This tests that the live ranges of the elements are updated correctly. 
    for i:=0; i<256; i++ {
        left_int_id := ctx.LoadInt(left_int)
        right_int_id := ctx.LoadInt(right_int)
        
        sum_el := ctx.Eval(SSA_ADD, left_int_id, right_int_id)
        
        if sum_el != i+2 {
            t.Errorf("Id of sum element is wrong, got: %v wanted: %v", sum_el, i+2)
        }
        
        if id:= ctx.Elements[left_int_id].LiveEnd; id != i+2 {
            t.Errorf("Live range of left int is wrong, got: %v wanted: %v", id, i+2)
        }
        
        if id:= ctx.Elements[right_int_id].LiveEnd; id != i+2 {
            t.Errorf("Live range of right int is wrong, got: %v wanted: %v", id, i+2)
        }        
    }    
}

func TestRegisterAllocation(t *testing.T) {    
    ctx := new (SsaContext)
    ctx.Init()
    
    some_int  := big.NewInt(1000)
    some_int_id := ctx.LoadInt(some_int)
            
    old_sum_el := 0
       
    // This creates a pathological chained expression that looks like: 1000 + 1000 + 1000 ... 256 times ... + 1000
    // which requires the allocator to activate and deactivate elements constantly, while still keeping one very
    // long lived element in a register. 
    for i:=0; i<256; i++ {        
        if old_sum_el == 0 {
            old_sum_el = ctx.Eval(SSA_ADD, some_int_id, some_int_id)
        } else {
            old_sum_el = ctx.Eval(SSA_ADD, some_int_id, old_sum_el)
        }       
    }
    
    // Really stress the allocator by allowing only 4 registers.
    // This seems to be the minimum necessary to solve this problem without
    // spilling registers.
    ctx.AllocateRegisters(4)  
}

