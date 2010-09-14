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

import (
        "big"
        "fmt"
)

type FloatObject struct {
    ObjectData
    Value float64 
}

// Convert float to int
func (o *FloatObject) AsInt() (*big.Int) {
    return big.NewInt(int64(o.Value))
}

// Convert float to float (identity transform)
func (o *FloatObject) AsFloat() (float64) {
    return o.Value
}

// Convert float to string
func (o *FloatObject) AsString() (string) {
    return fmt.Sprint(o.Value)
}

///////// Rich Comparison Interface ///////////

func (o *FloatObject) Lt(r Object) (bool) {
    return o.Value < r.AsFloat()
}

func (o *FloatObject) Gt(r Object) (bool) {
    return o.Value > r.AsFloat()
}

func (o *FloatObject) Eq(r Object) (bool) {
    return o.Value == r.AsFloat()
}

func (o *FloatObject) Neq(r Object) (bool) {
    return o.Value != r.AsFloat()
}

func (o *FloatObject) Lte(r Object) (bool) {
    return o.Value <= r.AsFloat()
}

func (o *FloatObject) Gte(r Object) (bool) {
    return o.Value >= r.AsFloat()
}

///////// Binary Arithmetic Interface ///////////

func (o *FloatObject) Add(r Object) (Object) {
    result := new (FloatObject)
    result.Value = o.Value + r.AsFloat()
    
    return result
}

func (o *FloatObject) Sub(r Object) (Object) {
    result := new (FloatObject)
    result.Value = o.Value - r.AsFloat()
    
    return result
}

func (o *FloatObject) Mul(r Object) (Object) {
    result := new (FloatObject)
    result.Value = o.Value * r.AsFloat()
    
    return result
}

func (o *FloatObject) Div(r Object) (Object) {
    result := new (FloatObject)
    result.Value = o.Value / r.AsFloat()
    
    return result
}

func (o *FloatObject) FloorDiv(r Object) (Object) {
    // Python says that the result of floor division
    // is always an integer.
    result := new (IntObject)
    result.Int = big.NewInt(int64(o.Value / r.AsFloat()))
    
    return result
}

func (o *FloatObject) Mod(r Object) (Object) {
    // We actually need to throw an exception, since
    // you can't mod two float objects.
    result := new (FloatObject)
    result.Value = 0
    
    return result
}


