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

   Contains the virtual machine itself, including the register defs.
  
*/

import "types"

type Machine struct {
    Register    [32]*types.Object     
    Pred        [32]bool
    
    NextInstruction int
}

func (m *Machine) Dispatch(c* CodeStream) {
    var op := c.ReadByte()
    
    switch(op) {
        case NOP:
        case LOAD: {
            var id          := c.ReadInt();
            var target_reg  := c.ReadByte();
            
            m.Register[target_reg] = c.Locals[id];
        }
        
        
            
    }
}
