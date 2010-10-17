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

   This file provides the implementation of the x86 in-memory assembler.
*/

package python

import "bytes"
import "encoding/binary"

type RegisterId uint8

const (
	x86_eax RegisterId = iota
	x86_ecx
	x86_edx
	x86_ebx
	x86_esp
	x86_ebp
	x86_esi
	x86_edi

	x64_r8
	x64_r9
	x64_r10
	x64_r11
	x64_r12
	x64_r13
	x64_r14
	x64_r15
)

const (
	vec_xmm0 = iota
	vec_xmm1
	vec_xmm2
	vec_xmm3
	vec_xmm4
	vec_xmm5
	vec_xmm6
	vec_xmm7
)


const (
	x86_conditionO = iota
	x86_conditionNO
	x86_conditionB
	x86_conditionAE
	x86_conditionE
	x86_conditionNE
	x86_conditionBE
	x86_conditionA
	x86_conditionS
	x86_conditionNS
	x86_conditionP
	x86_conditionNP
	x86_conditionL
	x86_conditionGE
	x86_conditionLE
	x86_conditionG

	x86_conditionC  = x86_conditionB
	x86_conditionNC = x86_conditionAE
)

type OneByteOpcodeId uint8

// One-byte opcodes
const (
	x86_ADD_EvGv                     OneByteOpcodeId = 0x01
	x86_ADD_GvEv                     = 0x03
	x86_OR_EvGv                      = 0x09
	x86_OR_GvEv                      = 0x0B
	x86_2BYTE_ESCAPE                 = 0x0F
	x86_AND_EvGv                     = 0x21
	x86_AND_GvEv                     = 0x23
	x86_SUB_EvGv                     = 0x29
	x86_SUB_GvEv                     = 0x2B
	x86_PRE_PREDICT_BRANCH_NOT_TAKEN = 0x2E
	x86_XOR_EvGv                     = 0x31
	x86_XOR_GvEv                     = 0x33
	x86_CMP_EvGv                     = 0x39
	x86_CMP_GvEv                     = 0x3B
	x64_PRE_REX                      = 0x40
	x86_PUSH_EAX                     = 0x50
	x86_Px86_EAX                     = 0x58
	x64_MOVSXD_GvEv                  = 0x63
	x86_PRE_OPERAND_SIZE             = 0x66
	x86_PRE_SSE_66                   = 0x66
	x86_PUSH_Iz                      = 0x68
	x86_IMUL_GvEvIz                  = 0x69
	x86_GROUP1_EbIb                  = 0x80
	x86_GROUP1_EvIz                  = 0x81
	x86_GROUP1_EvIb                  = 0x83
	x86_TEST_EvGv                    = 0x85
	x86_XCHG_EvGv                    = 0x87
	x86_MOV_EvGv                     = 0x89
	x86_MOV_GvEv                     = 0x8B
	x86_LEA                          = 0x8D
	x86_GROUP1A_Ev                   = 0x8F
	x86_CDQ                          = 0x99
	x86_MOV_EAXOv                    = 0xA1
	x86_MOV_OvEAX                    = 0xA3
	x86_MOV_EAXIv                    = 0xB8
	x86_GROUP2_EvIb                  = 0xC1
	x86_RET                          = 0xC3
	x86_GROUP11_EvIz                 = 0xC7
	x86_INT3                         = 0xCC
	x86_GROUP2_Ev1                   = 0xD1
	x86_GROUP2_EvCL                  = 0xD3
	x86_CALL_rel32                   = 0xE8
	x86_JMP_rel32                    = 0xE9
	x86_PRE_SSE_F2                   = 0xF2
	x86_HLT                          = 0xF4
	x86_GROUP3_EbIb                  = 0xF6
	x86_GROUP3_Ev                    = 0xF7
	x86_GROUP3_EvIz                  = 0xF7 // x86_GROUP3_Ev has an immediate when instruction is a test. 
	x86_GROUP5_Ev                    = 0xFF
)

type TwoByteOpcodeId uint8

// Two-byte op codes
const (
	x86_MOVSD_VsdWsd    TwoByteOpcodeId = 0x10
	x86_MOVSD_WsdVsd    = 0x11
	x86_CVTSI2SD_VsdEd  = 0x2A
	x86_CVTTSD2SI_GdWsd = 0x2C
	x86_UCOMISD_VsdWsd  = 0x2E
	x86_ADDSD_VsdWsd    = 0x58
	x86_MULSD_VsdWsd    = 0x59
	x86_SUBSD_VsdWsd    = 0x5C
	x86_DIVSD_VsdWsd    = 0x5E
	x86_SQRTSD_VsdWsd   = 0x51
	x86_XORPD_VpdWpd    = 0x57
	x86_MOVD_VdEd       = 0x6E
	x86_MOVD_EdVd       = 0x7E
	x86_JCC_rel32       = 0x80
	x86_SETCC           = 0x90
	x86_IMUL_GvEv       = 0xAF
	x86_MOVZX_GvEb      = 0xB6
	x86_MOVZX_GvEw      = 0xB7
	x86_PEXTRW_GdUdIb   = 0xC5
)

func jccRel32(cond uint8) TwoByteOpcodeId {
	return (TwoByteOpcodeId)(x86_JCC_rel32 + cond)
}

func setccOpcode(cond uint8) TwoByteOpcodeId {
	return (TwoByteOpcodeId)(x86_SETCC + cond)
}

type GroupOpcodeId uint8

const (
	x86_GROUP1_OP_ADD = 0
	x86_GROUP1_OP_OR  = 1
	x86_GROUP1_OP_ADC = 2
	x86_GROUP1_OP_AND = 4
	x86_GROUP1_OP_SUB = 5
	x86_GROUP1_OP_XOR = 6
	x86_GROUP1_OP_CMP = 7

	x86_GROUP1A_OP_POP = 0

	x86_GROUP2_OP_SHL = 4
	x86_GROUP2_OP_SHR = 5
	x86_GROUP2_OP_SAR = 7

	x86_GROUP3_OP_TEST = 0
	x86_GROUP3_OP_NOT  = 2
	x86_GROUP3_OP_NEG  = 3
	x86_GROUP3_OP_IDIV = 7

	x86_GROUP5_OP_CALLN = 2
	x86_GROUP5_OP_JMPN  = 4
	x86_GROUP5_OP_PUSH  = 6

	x86_GROUP11_MOV = 0
)

const (
    ModRmMemoryNoDisp = iota
    ModRmMemoryDisp8
    ModRmMemoryDisp32
    ModRmRegister
)

const ( noBase  RegisterId = x86_ebp
        hasSib             = x86_esp
        noIndex            = x86_esp

        noBase2            = x64_r13
        hasSib2            = x64_r12
)

/*******************************************************************
 * Instruction buffer 
 *******************************************************************/
type X86Buffer struct {
    *bytes.Buffer
    
    IsX64   bool
}

/*******************************************************************
 * Instruction formatting structures
 *******************************************************************/

type JmpSrc struct {
	offset int
}

type JmpDst struct {
	offset int
	used   bool
}

/*******************************************************************
 * Instruction formatting functions
 *******************************************************************/

// Byte-operands:
//
// These methods format byte operations.  Byte operations differ from the normal
// formatters in the circumstances under which they will decide to emit REX prefixes.
// These should be used where any register operand signifies a byte register.
//
// The disctinction is due to the handling of register numbers in the range 4..7 on
// x86-64.  These register numbers may either represent the second byte of the first
// four registers (ah..bh) or the first byte of the second four registers (spl..dil).
//
// Since ah..bh cannot be used in all permutations of operands (specifically cannot
// be accessed where a REX prefix is present), these are likely best treated as
// deprecated.  In order to ensure the correct registers spl..dil are selected a
// REX prefix will be emitted for any byte register operand in the range 4..15.
//
// These formatters may be used in instructions where a mix of operand sizes, in which
// case an unnecessary REX will be emitted, for example:
//     movzbl %al, %edi
// In this case a REX will be planted since edi is 7 (and were this a byte operand
// a REX would be required to specify dil instead of bh).  Unneeded REX prefixes will
// be silently ignored by the processor.
//
// Address operands should still be checked using regRequiresRex(), while byteRegRequiresRex()
// is provided to check byte register operands.

func (buf *X86Buffer) fmtOp8(opcode OneByteOpcodeId, groupOp GroupOpcodeId, rm RegisterId) {    
    buf.emitRexIf(buf.byteRegRequiresRex(rm), 0, 0, rm)
    buf.WriteByte(byte(opcode))
    buf.registerModRM(RegisterId(groupOp), rm)
}

func (buf *X86Buffer) fmtExtOp8(opcode TwoByteOpcodeId, reg RegisterId, rm RegisterId) {    
    buf.emitRexIf(buf.byteRegRequiresRex(reg)||buf.byteRegRequiresRex(rm), reg, 0, rm)
    buf.WriteByte(x86_2BYTE_ESCAPE)
    buf.WriteByte(byte(opcode))
    buf.registerModRM(reg, rm)
}

func (buf *X86Buffer) fmtExtGroupOp8(opcode TwoByteOpcodeId, groupOp GroupOpcodeId, rm RegisterId) {    
    buf.emitRexIf(buf.byteRegRequiresRex(rm), 0, 0, rm)
    buf.WriteByte(x86_2BYTE_ESCAPE)
    buf.WriteByte(byte(opcode))
    buf.registerModRM(RegisterId(groupOp), rm)
}

func (buf *X86Buffer) putModRm(mode, reg, rm RegisterId) {
    buf.WriteByte(uint8((int(mode) << 6) | ((int(reg) & 7) << 3) | (int(rm) & 7)))
}

func (buf *X86Buffer) putModRmSib(mode, reg, base, index RegisterId, scale int) {            
    buf.putModRm(mode, reg, hasSib)
    buf.WriteByte(uint8((int(scale) << 6) | ((int(index) & 7) << 3) | (int(base) & 7)))
}

func (buf *X86Buffer) registerModRM(reg, rm RegisterId) {
    buf.putModRm(ModRmRegister, reg, rm)
}

// Immediates:
//
// An immediate should be appended where appropriate after an op has been emitted.
// The writes are unchecked since the opcode formatters above will have ensured space.

func immediate(buf *bytes.Buffer, imm int8) {
   binary.Write(buf, binary.LittleEndian, imm)     
}

func immediate16(buf *bytes.Buffer, imm int16) {
    binary.Write(buf, binary.LittleEndian, imm)    
}

func immediate32(buf *bytes.Buffer, imm int32) {
    binary.Write(buf, binary.LittleEndian, imm)    
}

func immediate64(buf *bytes.Buffer, imm int64) {
    binary.Write(buf, binary.LittleEndian, imm)    
}

func immediateRel32(buf *bytes.Buffer) JmpSrc {
    binary.Write(buf, binary.LittleEndian, 0)
    return JmpSrc { buf.Len() }
}

// Registers r8 & above require a REX prefixe.
func (buf *X86Buffer) regRequiresRex(reg RegisterId) bool {
    if buf.IsX64 {
        return (reg >= x64_r8)
    }
    
    return false
}

// Byte operand register spl & above require a REX prefix (to prevent the 'H' registers be accessed).
func (buf *X86Buffer) byteRegRequiresRex(reg RegisterId) bool {
    if buf.IsX64 {
        return (reg >= x86_esp)
    }
    
    return false
}

// Format a REX prefix byte.
func (buf *X86Buffer) emitRex(w bool, r, x, b RegisterId) {
    v := x64_PRE_REX | ((int(r)>>3)<<2) | ((int(x)>>3)<<1) | (int(b)>>3)
    if w {
        v |= 1<<3
    } 
    
    buf.WriteByte(uint8(v));
}

// Used to plant a REX byte with REX.w set (for 64-bit operations).
func (buf *X86Buffer) emitRexW(r, x, b RegisterId) {
    buf.emitRex(true, r, x, b);
}

// Used for operations with byte operands - use byteRegRequiresRex() to check register operands,
// regRequiresRex() to check other registers (i.e. address base & index).
func (buf *X86Buffer) emitRexIf(condition bool, r, x, b RegisterId) {
    if condition {
        buf.emitRex(false, r, x, b)
    }
}

// Used for word sized operations, will plant a REX prefix if necessary (if any register is r8 or above).
func (buf *X86Buffer) emitRexIfNeeded(r, x, b RegisterId) {
    buf.emitRexIf(buf.regRequiresRex(r) || buf.regRequiresRex(x) || buf.regRequiresRex(b), r, x, b);
}
