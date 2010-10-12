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

   This module implements SSA operations.  We trade memory for speed in the
   SSA representation, and we make an optimization pass to reduce the code
   before register allocation (as part of the bytecode translation pass.)  
*/

package python

import (
	"big"
	"container/vector"
	"fmt"
)

const (
	SSA_CALL = iota
	SSA_SPILL
	SSA_FILL
	SSA_LOAD
	SSA_STORE
	SSA_ALU_MARK
	SSA_ADD
	SSA_SUB
	SSA_MUL
	SSA_DIV
	SSA_MOD
	SSA_POW
	SSA_AND
	SSA_OR
	SSA_XOR
	SSA_NOT
	SSA_GET
	SSA_SET
	SSA_IDX
)

const (
	SSA_TYPE_ELEMENT = iota
	SSA_TYPE_CLASS
	SSA_TYPE_INTEGER
	SSA_TYPE_STRING
	SSA_TYPE_BUFFER
	SSA_TYPE_FLOAT
	SSA_TYPE_COMPLEX
	SSA_TYPE_BOOL
	SSA_TYPE_NONE
	SSA_TYPE_UNKNOWN
)

// The SsaElement is a single assignment, which may include
// a single operation.  The element represents the result of
// the operation.  The simplest operation is just "SSA_ASSIGN"
// which causes this element to take on the value of the src1
// operand.  All other elements involve both operands, and the
// results of some operation on them. 
type SsaElement struct {
	// The operation to perform, one of SSA_XXX
	Op uint

	// The two source operands
	Src1, Src2 int

	// The type of the source operands, one of SSA_TYPE_XXX
	Src1Type, Src2Type uint

	// Flags set if this element is ever read, and if it is known to be
	// constant at compile time.  By definition an element is always written to,
	// since an SSA element will never be created without a write.
	// Pinned means that the instruction will always be emitted (never optimized
	// away.)
	WasRead, IsConst, Pinned bool

	// These indicate at what point this element becomes live (is first initialized)
	// and when it dies (is never used again.)  These are important values to know
	// so that we can maintain the active list during register allocation.  The value
	// used is the index of the first SSA where this element is used, and the last index
	// where this element is used.
	LiveStart, LiveEnd int

	ActiveStart, ActiveEnd int

	// The registers allocated to this element. 0 means unallocated, since only 0 values can
	// be mapped to register 0.  A single element may be spilled, meaning that it is later
	// mapped back in as a _source_ to different registers.  
	DstRegister, Src1Register, Src2Register int

	// The address of this element in the current code stream
	Address int
}

// Helps to track items which had to be spilled
// from the register bank during register allocation.
type SsaMapContext struct {

	// Storage for the free spill slots 
	FreeSpillSlots *vector.IntVector

	// At any given time, some elements
	// must not be spilled because they
	// are needed by the current instruction
	NoSpillElements map[int]bool

	// Map of spill slots to SSA element
	SpillMap map[int]int

	// Tracks old_ssa_id -> new_ssa_id values so
	// we can rename the parameters correctly during rewrite.
	RenameMap map[int]int

	// The list of free regs is kept here
	FreeRegs *vector.IntVector

	// Store the active SSA elements in this list.
	ActiveElements *vector.Vector
}


func (s *SsaMapContext) Init() {
	s.FreeSpillSlots = new(vector.IntVector)
	s.FreeRegs = new(vector.IntVector)
	s.ActiveElements = new(vector.Vector)

	s.NoSpillElements = make(map[int]bool, 8)
	s.SpillMap = make(map[int]int, 8)
	s.RenameMap = make(map[int]int, 8)
}

type SsaContext struct {
	LastElementId int
	Elements      []*SsaElement
	Ints          *vector.Vector
	Floats        *vector.Vector
	Strings       *vector.StringVector
	Names         *vector.StringVector

	// The maps below are actually maps from
	// the values to the SsaElements created
	// to load them into an SSA "register".


	IntIdx    map[*big.Int]int
	FloatIdx  map[float64]int
	StringIdx map[string]int
	NameIdx   map[string]int

	// How many slots are needed for some
	// code object in order to spill
	SpillRoomNeeded int

	// This is set when the live checks performed
	// in Write should be turned off.  This is
	// useful during register allocation and optimization.
	DisableLiveCheck bool
}

func (ctx *SsaContext) Init() {
	ctx.Elements = make([]*SsaElement, 128, 128)
	ctx.Ints = new(vector.Vector)
	ctx.Floats = new(vector.Vector)
	ctx.Strings = new(vector.StringVector)
	ctx.Names = new(vector.StringVector)

	ctx.IntIdx = make(map[*big.Int]int, 16)
	ctx.FloatIdx = make(map[float64]int, 16)
	ctx.StringIdx = make(map[string]int, 16)
	ctx.NameIdx = make(map[string]int, 16)
}

func (ctx *SsaContext) Write(el *SsaElement) int {
	// Grow the element slice if we are out of space
	if ctx.LastElementId >= len(ctx.Elements) {
		tmp := make([]*SsaElement, ctx.LastElementId+128, ctx.LastElementId+128)

		for i := 0; i < ctx.LastElementId; i++ {
			tmp[i] = ctx.Elements[i]
		}

		ctx.Elements = tmp
	}

	if !ctx.DisableLiveCheck {
		// Initialize the live ranges
		el.LiveStart = ctx.LastElementId
		el.LiveEnd = ctx.LastElementId

		// Update the element(s) that this element references as having been read, and
		// update their live range too.
		if el.Op > SSA_ALU_MARK {
			if el.Src1Type == SSA_TYPE_ELEMENT {
				ctx.Elements[el.Src1].WasRead = true
				ctx.Elements[el.Src1].LiveEnd = ctx.LastElementId
			}
			if el.Src2Type == SSA_TYPE_ELEMENT {
				ctx.Elements[el.Src2].WasRead = true
				ctx.Elements[el.Src2].LiveEnd = ctx.LastElementId
			}
		}
	}

	// Write a new element    
	el.Address = ctx.LastElementId
	ctx.Elements[ctx.LastElementId] = el
	ctx.LastElementId++

	return el.Address
}

func (ctx *SsaContext) Eval(op uint, src1, src2 int) int {

	el := new(SsaElement)

	el.Op = op
	el.Src1 = src1
	el.Src2 = src2

	// All ALU/FPU operations' operands are elements.  Only
	// LOAD/STORE deals with other types.
	el.Src1Type = SSA_TYPE_ELEMENT
	el.Src2Type = SSA_TYPE_ELEMENT

	return ctx.Write(el)
}

func (ctx *SsaContext) Spill(to_slot, from_register int) int {

	el := new(SsaElement)

	el.Op = SSA_SPILL
	el.Src1 = to_slot
	el.DstRegister = from_register

	return ctx.Write(el)
}

func (ctx *SsaContext) Fill(from_slot, to_register int) int {

	el := new(SsaElement)

	el.Op = SSA_FILL
	el.Src1 = from_slot
	el.DstRegister = to_register

	return ctx.Write(el)
}

func (ctx *SsaContext) LoadInt(v *big.Int) int {
	idx, present := ctx.IntIdx[v]

	if !present {
		// Save the integer in the array so we know what the actual
		// value should be        
		idx = len(ctx.IntIdx)
		ctx.Ints.Push(v)

		// Create a new SSA element to store the actual action of 
		// loading a literal int
		el := new(SsaElement)

		el.Op = SSA_LOAD
		el.Src1 = idx
		el.Src1Type = SSA_TYPE_INTEGER

		// Map the new element to the value    
		idx = ctx.Write(el)
		ctx.IntIdx[v] = idx
	}

	return idx
}

// Generates a spill instruction.  Decides what to spill, and generates an instruction to save
// the spilled value.  The return value is the newly freed register.  
func (ctx *SsaContext) generateSpill(mc *SsaMapContext) int {

	// Find a register to spill.  Our heuristic is to
	// choose the register with the longest lifetime. That
	// seems counter-intuitive, but http://www.cs.ucla.edu/~palsberg/course/cs132/linearscan.pdf
	// indicates that it performs best.  Assuming I understood the
	// paper, of course.
	var spill_el *SsaElement = nil
	spilled_el_index := 0

	for i := 0; i < mc.ActiveElements.Len(); i++ {
		candidate_el := mc.ActiveElements.At(i).(*SsaElement)

		if _, present := mc.NoSpillElements[candidate_el.Address]; present {
			// If we don't have an element to spill yet, or if the current
			// element is a better candidate, choose it.
			if spill_el == nil || spill_el.LiveEnd < candidate_el.LiveEnd {
				spill_el = candidate_el
				spilled_el_index = i
			}
		}
	}
	
	if spill_el == nil {
	   panic("There are no spillable registers.")
	}

	free_slot := 0

	// Once we've chose a register, we need to figure out where to spill the
	// data to.  We try to make this reasonably optimal, but we can grow the
	// spill area as needed.  (Something not true about our register set. :-D)
	if mc.FreeSpillSlots.Len() == 0 {
		// No free spill slots, grow it.
		free_slot = len(mc.SpillMap)
	} else {
		free_slot = mc.FreeSpillSlots.Pop()
	}

	mc.SpillMap[spill_el.Address] = free_slot

	// Now emit a spill instruction
	// so that we don't lose the work done.            
	ctx.Spill(free_slot, spill_el.DstRegister)

	// Make sure to track how much spill room is needed
	if ctx.SpillRoomNeeded < len(mc.SpillMap) {
		ctx.SpillRoomNeeded = len(mc.SpillMap)
	}

	// Remove it from the active list
	mc.ActiveElements.Delete(spilled_el_index)

	fmt.Printf("spilled: %v\n", spill_el.Address)

	// Return the newly freed register number    
	return spill_el.DstRegister
}

// Generates a fill instruction.  Previously the value must have been spilled out to the save area.  An
// instruction is emitted to load it back into the register set.  Other registers may be spilled in order
// to bring the spilled value back in.  Returns the id of the element that generated the fill.  This id
// should be used as the new source value of an SsaElement that depends on the spilled value.
func (ctx *SsaContext) generateFill(el *SsaElement, mc *SsaMapContext) int {

	// Figure out where the element was 
	// spilled to.
	free_slot := mc.SpillMap[el.Address]
	mc.FreeSpillSlots.Push(free_slot)

	target_reg := 0

	// Find a free register (possibly by spilling another register.)
	if mc.FreeRegs.Len() == 0 {
		target_reg = ctx.generateSpill(mc)
	} else {
		target_reg = mc.FreeRegs.Pop()
	}

	// Remove the element from the map
	mc.SpillMap[el.Address] = 0, false

	// Activate the element.
	mc.ActiveElements.Push(el)

	fmt.Printf("filled: %v\n", el.Address)

	// Write the fill instruction
	return ctx.Fill(free_slot, target_reg)
}


// Performs a linear-scan allocation of registers.  Only one pass is used to allocate registers to all
// SSA instructions.
func (ctx *SsaContext) AllocateRegisters(num_regs int) *SsaContext {

	// We create a new context so that we can rewrite the SSA stream into it.  This is because
	// we expect that we will need to spill at least one SSA into a temporary space.  A possible
	// future optimization of this code would be to have the Strahler number calculated by the
	// AST traversal phase so we know if we will need to spill or not.  Of course, we also take
	// this opportunity to do some optimizations that require rewriting the stream anyway (like 
	// dead code elimination.)

	new_ctx := new(SsaContext)
	new_ctx.Init()
	new_ctx.DisableLiveCheck = true

	// The list of spilled elements is kept here
	mc := new(SsaMapContext)
	mc.Init()

	// Push all the registers except 0 onto the free list. We assume the 0 register
	// is reserved for the 0 value, thus it is never available.
	for i := 1; i < num_regs; i++ {
		mc.FreeRegs.Push(i)
	}

	for ssa_id := 0; ssa_id < ctx.LastElementId; ssa_id++ {
		old_el := ctx.Elements[ssa_id]

		// First, check to see if this element is ever read.
		if !old_el.Pinned && !old_el.WasRead {
			// This element was never looked at, so we can
			// skip it.
			continue
		}

		// Create a new element to copy the
		// old one into
		el := new(SsaElement)
		*el = *old_el

		///////////////////

		new_active_elements := new(vector.Vector)

		// First remove any elements whose LiveEnd value is less than the 
		// current ssa_id index
		for i := 0; i < mc.ActiveElements.Len(); i++ {

			candidate_el := mc.ActiveElements.At(i).(*SsaElement)

			fmt.Printf("%v: live: %v,%v\n", ssa_id, candidate_el.LiveStart, candidate_el.LiveEnd)

			if candidate_el.LiveEnd >= ssa_id {
				new_active_elements.Push(candidate_el)
			} else {
				// Indicate that this register is free again
				mc.FreeRegs.Push(candidate_el.DstRegister)
				el.ActiveEnd = ssa_id
			}
		}

		// Use the new list as our active elements list
		mc.ActiveElements = new_active_elements

		// Update the active start address
		el.ActiveStart = ssa_id

		// Process any renames and fills
		if el.Op > SSA_ALU_MARK {
			// Check for (and perform) any needed renames.
			if new_src1_name, present := mc.RenameMap[el.Src1]; present {
				el.Src1 = new_src1_name
			}

			if new_src2_name, present := mc.RenameMap[el.Src2]; present {
				el.Src2 = new_src2_name
			}

			mc.NoSpillElements[el.Src1] = true
			mc.NoSpillElements[el.Src2] = true

			// Check to see if we need to fill some registers from the
			// spill area in order to process this instruction.  If so, 
			// we _may_ need to spill one or two registers in order to
			// have the space we need to fill for this instruction.	        
			if _, spilled := mc.SpillMap[el.Src1]; spilled {
				el.Src1 = new_ctx.generateFill(new_ctx.Elements[el.Src1], mc)
			}

			if _, spilled := mc.SpillMap[el.Src2]; spilled {
				el.Src2 = new_ctx.generateFill(new_ctx.Elements[el.Src2], mc)
			}
			
			//// PROBLEM:
			// Some instruction results are swapped in and out of the register
			// file.  This means that at certain points they have been moved
			// to different registers.  We need to keep track of the fact that
			// the register is different during various intervals.  We need to
			// know WHAT it is and WHEN it is that value.
			// This means that we need to put the filled values into the activated
			// records and rename the first src to the filled src.  The filled src's
			// live range should be set to end at the same place as the original source
			// to maintain the heuristic.  If we spill the filled record, then we should 
			// delete the rename map too.  However, the original source may have been
			// renamed due to optimizations or dead-code elimination.  So somehow we need
			// to get back to the original rename.
			
		}

		// Figure out what register this instruction should go into
		if mc.FreeRegs.Len() == 0 {
			el.DstRegister = new_ctx.generateSpill(mc)
		} else {
			el.DstRegister = mc.FreeRegs.Pop()
		}

		// Track the register in the new and old context.
		old_el.DstRegister = el.DstRegister

		// Write the possibly renamed element into the new context                
		mc.RenameMap[ssa_id] = new_ctx.Write(el)

		// Push the current eement into the active elements list.
		// Do this here so that it does not get considered for 
		// spilling.
		mc.ActiveElements.Push(el)

		// Clear out the no-spill list.
		mc.NoSpillElements[el.Src1] = false, false
		mc.NoSpillElements[el.Src2] = false, false
	}

	return new_ctx
}
