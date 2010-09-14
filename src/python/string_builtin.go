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

   This file provides the implementation of the string built-in object
   type.
*/

package python

import (
        "big"
        "fmt"
)

type StringObject struct {
    ObjectData
    Value string 
}

func NewString(value string) (*StringObject) {
    str := new(StringObject)
    str.ObjectData.Init()
    str.Value = value
    
    return str
}

// Convert string to int
func (o *StringObject) AsInt() (*big.Int) {
    value := big.NewInt(0)
    value.SetString(o.Value, 0)
    
    return value
}

// Convert string to float
func (o *StringObject) AsFloat() (float64) {
    var value float64
    
    fmt.Scan(o.Value, value)
    return value
}

// Convert string to string (identity transform)
func (o *StringObject) AsString() (string) {
    return o.Value
}

///////// Rich Comparison Interface ///////////

func (o *StringObject) Lt(r Object) (bool) {
    return o.Value < r.AsString()
}

func (o *StringObject) Gt(r Object) (bool) {
    return o.Value > r.AsString()
}

func (o *StringObject) Eq(r Object) (bool) {
    return o.Value == r.AsString()
}

func (o *StringObject) Neq(r Object) (bool) {
    return o.Value != r.AsString()
}

func (o *StringObject) Lte(r Object) (bool) {
    return o.Value <= r.AsString()
}

func (o *StringObject) Gte(r Object) (bool) {
    return o.Value >= r.AsString()
}

///////// Binary Arithmetic Interface ///////////

func (o *StringObject) Add(r Object) (Object) {    
    return NewString(o.Value + r.AsString())
}

func (o *StringObject) Sub(r Object) (Object) {    
    return NewString("")
}

func (o *StringObject) Mul(r Object) (Object) { 
    result := ""
    reps   := r.AsInt().Int64()
    
    for i:=int64(0); i < reps; i+=1 {
        result+=o.Value
    } 
    return NewString(result)
}

func (o *StringObject) Div(r Object) (Object) {
    return NewString("")
}

func (o *StringObject) FloorDiv(r Object) (Object) {
    return NewString("")
}

func (o *StringObject) Mod(r Object) (Object) {
    return NewString(o.Value)
}


