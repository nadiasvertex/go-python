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

   This file provides the implementation of the integer built-in object
   type.
*/

package python

import "big"

type IntObject struct {
    ObjectData
    *big.Int 
}

func NewIntObject() (*IntObject) {
    r := new (IntObject)
    r.Int = big.NewInt(0)
    
    return r
}

// Convert int to int (identity transform)
func (o *IntObject) AsInt() (*big.Int) {
    return o.Int
}

// Convert int to float
func (o *IntObject) AsFloat() (float64) {
    return float64(o.Int64())
}

// Convert int to string
func (o *IntObject) AsString() (string) {
    return o.String()
}


///////// Rich Comparison Interface ///////////

func (o *IntObject) Lt(r Object) (bool) {
    return o.Cmp(r.AsInt()) == -1
}

func (o *IntObject) Gt(r Object) (bool) {
    return o.Cmp(r.AsInt()) == 1
}

func (o *IntObject) Eq(r Object) (bool) {
    return o.Cmp(r.AsInt()) == 0
}

func (o *IntObject) Neq(r Object) (bool) {
    return o.Cmp(r.AsInt()) != 0
}

func (o *IntObject) Lte(r Object) (bool) {
    return o.Cmp(r.AsInt()) <= 0
}

func (o *IntObject) Gte(r Object) (bool) {
    return o.Cmp(r.AsInt()) >= 0
}

///////// Binary Arithmetic Interface ///////////

func (o *IntObject) Add(r Object) (Object) {
    result := NewIntObject()
    result.Int.Add(o.Int, r.AsInt())
    
    return result
}

func (o *IntObject) Sub(r Object) (Object) {
    result := NewIntObject()
    result.Int.Sub(o.Int, r.AsInt())
    
    return result
}

func (o *IntObject) Mul(r Object) (Object) {
    result := NewIntObject()
    result.Int.Mul(o.Int, r.AsInt())
    
    return result
}

func (o *IntObject) Div(r Object) (Object) {
    // Python says that the result of a '/' operation
    // is always a FloatObject, irregardless of whether
    // the input is an integer or float
    result := new (FloatObject)
    result.Value = float64(o.Int.Int64()) / r.AsFloat()
    
    return result
}

func (o *IntObject) FloorDiv(r Object) (Object) {
    // This is the // operation, which results in an 
    // integer.
    result := NewIntObject()    
    result.Int.Div(o.Int, r.AsInt())
    
    return result
}

func (o *IntObject) Mod(r Object) (Object) {
    result := NewIntObject()
    result.Int.Mod(o.Int, r.AsInt())
    
    return result
}



