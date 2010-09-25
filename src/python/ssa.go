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

import "big"

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
    Src1, Src2  uint
    
    // The type of the source operands, one of SSA_TYPE_XXX
    Src1Type, Src2Type  uint 
    
    // Flags set if this element is ever read, and if it is known to be
    // constant at compile time.  By definition an element is written to,
    // since an SSA element will never be created without a write.
    WasRead, IsConst    bool 
    
    // The register allocated to this element. 0 means unallocated, since only 0 values can
    // be mapped to register 0.  
    Register    uint
}

type SsaContext struct {
    LastElementId   int
    Elements        []*SsaElement
    Ints            []*big.Int
    Floats          []float64
    Strings         []string
    Names           []string
    
    
    IntIdx          map[*big.Int]int
    FloatIdx        map[float64]int
    StringIdx       map[string]int        
    NameIdx         map[string]int
}

func (ctx *SsaContext) Init() {
    ctx.Elements = make([]*SsaElement, 128, 128)
    ctx.Ints     = make([]*big.Int, 16, 16)
    ctx.Floats   = make([]float64, 16, 16)
    ctx.Strings  = make([]string, 16, 16)
    ctx.Names    = make([]string, 16, 16)
    
     
    ctx.IntIdx      = make(map[*big.Int]int, 16)    
    ctx.FloatIdx    = make(map[float64]int, 16)    
    ctx.StringIdx   = make(map[string]int, 16)     
    ctx.NameIdx     = make(map[string]int, 16)
}

func (ctx *SsaContext) Write(el *SsaElement) int {
    // Grow the element slice if we are out of space
    if ctx.LastElementId >= len(ctx.Elements) {
        tmp := make([]*SsaElement, ctx.LastElementId + 128, ctx.LastElementId + 128)
        
        for i:=0; i<ctx.LastElementId; i++ {
            tmp[i] = ctx.Elements[i]
        } 
        
        ctx.Elements = tmp
    }    
    
    // Write a new element
    ctx.Elements[ctx.LastElementId] = el
    ctx.LastElementId++
    
    return ctx.LastElementId-1
}

