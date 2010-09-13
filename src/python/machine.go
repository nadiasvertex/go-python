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

package python

import "encoding/binary"

// All instruction types
const instruction_mask  uint32 = 0x000003f
const pred_execute_mask uint32 = 0x0000040
const pred_reg_mask     uint32 = 0x0000F80

const pred_execute_shift    uint32 = 6
const pred_reg_shift        uint32 = 7

// Register mode instruction types
const source_reg1_mask  uint32 = 0x000F000
const source_reg2_mask  uint32 = 0x00F0000
const target_reg_mask   uint32 = 0x0F00000

const source_reg1_shift uint32 = 12
const source_reg2_shift uint32 = 16
const target_reg_shift  uint32 = 20


// Immediate mode instruction types
const imm_target_reg_mask   uint32 = 0x0000F000
const immediate_val_mask    uint32 = 0xFFFF0000

const imm_target_reg_shift  uint32 = 12
const immediate_val_shift   uint32 = 16


type Machine struct {
    Register    [16]*Object     
    Pred        [32]bool
    
    NextInstruction uint32
}

func (m *Machine) Dispatch(c* CodeStream) {
    var instruction uint32     
    binary.Read(c, binary.LittleEndian, &instruction);
    
    op := instruction & instruction_mask;
    
    var /*reg1, reg2,*/ reg3 uint32  
    var imm              uint16     
    
    // Decoder stage - decodes the instruction based on our instruction formats.
    switch {
        case op <=15:
        
        case op <=31:
            reg3 = (instruction & imm_target_reg_mask)>>imm_target_reg_shift
            imm  = uint16((instruction & immediate_val_mask)>>immediate_val_shift)
            
        default:
            //reg1 = (instruction & source_reg1_mask)>>source_reg1_shift
            //reg2 = (instruction & source_reg2_mask)>>source_reg2_shift
            reg3 = (instruction & target_reg_mask)>>target_reg_shift
    }
    
    // Execution stage - actually processes the instructions.
    switch op {
        case NOP:
        case LOAD:
            m.Register[reg3] = c.Locals[imm]            
        case BIND:
            c.Locals[imm] = m.Register[reg3]       
            
    }
}
