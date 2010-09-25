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
)

const (
    SSA_LOAD = iota
    SSA_STORE
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
    SSA_CALL   
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
    Op          uint 
    
    // The two source operands
    Src1, Src2  int
    
    // The type of the source operands, one of SSA_TYPE_XXX
    Src1Type, Src2Type  uint 
    
    // Flags set if this element is ever read, and if it is known to be
    // constant at compile time.  By definition an element is written to,
    // since an SSA element will never be created without a write.
    WasRead, IsConst    bool 
    
    // These indicate at what point this element becomes live (is first initialized)
    // and when it dies (is never used again.)  These are important values to know
    // so that we can maintain the active list during register allocation.  The value
    // used is the index of the first SSA where this element is used, and the last index
    // where this element is used.
    LiveStart, LiveEnd  int
    
    // The register allocated to this element. 0 means unallocated, since only 0 values can
    // be mapped to register 0.  
    Register    uint
}

type SsaContext struct {
    LastElementId   int
    Elements        []*SsaElement
    Ints            *vector.Vector
    Floats          *vector.Vector
    Strings         *vector.StringVector
    Names           *vector.StringVector
    
    // The maps below are actually maps from
    // the values to the SsaElements created
    // to load them into an SSA "register".
    
    
    IntIdx          map[*big.Int]int
    FloatIdx        map[float64]int
    StringIdx       map[string]int        
    NameIdx         map[string]int
}

func (ctx *SsaContext) Init() {
    ctx.Elements = make([]*SsaElement, 128, 128)
    ctx.Ints     = new(vector.Vector)
    ctx.Floats   = new(vector.Vector)
    ctx.Strings  = new(vector.StringVector)
    ctx.Names    = new(vector.StringVector)    
     
    ctx.IntIdx      = make(map[*big.Int]int, 16)    
    ctx.FloatIdx    = make(map[float64]int, 16)    
    ctx.StringIdx   = make(map[string]int, 16)     
    ctx.NameIdx     = make(map[string]int, 16)
}

func (ctx *SsaContext) Write(el *SsaElement) (el_id int) {
    // Grow the element slice if we are out of space
    if ctx.LastElementId >= len(ctx.Elements) {
        tmp := make([]*SsaElement, ctx.LastElementId + 128, ctx.LastElementId + 128)
        
        for i:=0; i<ctx.LastElementId; i++ {
            tmp[i] = ctx.Elements[i]
        } 
        
        ctx.Elements = tmp
    }    
    
    // Initialize the live ranges
    el.LiveStart = ctx.LastElementId
    el.LiveEnd = ctx.LastElementId 
            
    // Update the element(s) that this element references as having been read, and
    // update their live range too.
    if el.Op > SSA_STORE {
        if el.Src1Type ==  SSA_TYPE_ELEMENT {
            ctx.Elements[el.Src1].WasRead = true
            ctx.Elements[el.Src1].LiveEnd = ctx.LastElementId
        }
        if el.Src2Type ==  SSA_TYPE_ELEMENT {
            ctx.Elements[el.Src2].WasRead = true
            ctx.Elements[el.Src2].LiveEnd = ctx.LastElementId
        }
    }
    
    // Write a new element
    el_id = ctx.LastElementId
    ctx.Elements[ctx.LastElementId] = el
    ctx.LastElementId++
    
    return
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
        el := new (SsaElement)
    
	    el.Op       = SSA_LOAD
	    el.Src1     = idx
	    el.Src1Type = SSA_TYPE_INTEGER        
    
        // Map the new element to the value    
        idx           = ctx.Write(el)
        ctx.IntIdx[v] = idx      
    }   
    
    return idx
}
