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

   The types package implements all the different types of objects that
   are built-in to the Python language.
*/

package python

type ObjectData struct {
    Name string
    Attrs map[string]Object 
}

//  Object attribute getting interface.
type Getter interface {
    GetAttr(name string) (value Object, present bool)     
} 

// Object attribute setting interface.
type Setter interface {
    SetAttr(name string, value Object)     
}

// Object rich comparison interface
type RichComparer interface {
    Lt(r Object) (bool)
    Gt(r Object) (bool)
    Eq(r Object) (bool)
    Neq(r Object) (bool)
    Lte(r Object) (bool)
    Gte(r Object) (bool)
}

type BinaryArithmetic interface {
    Add(r Object) (Object)
    Sub(r Object) (Object)
    Mul(r Object) (Object)
    Div(r Object) (Object)
    Mod(r Object) (Object)
}

type Converter interface {
    AsInt() (int)    
}

// Object composite interface
type Object interface {
    Getter
    Setter
    RichComparer
    BinaryArithmetic
    Converter
}

// Get the value of an object's attribute.
func (o *ObjectData) GetAttr(name string) (value Object, present bool) {
    value, present = o.Attrs[name]
    return  
}

// Set the value of an object's attribute.
func (o *ObjectData) SetAttr(name string, value Object) {
    o.Attrs[name] = value
    return  
}

/*
// Lookup the less than operator and execute it, if one exists.
func (o *Object) Lt(l, r *Object) (bool) {
    if cmp, present := o.GetAttr("__lt__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return false//l  < r
}

// Lookup the greater than operator and execute it, if one exists.
func (o *Object) Gt(l, r *Object) (bool) {
    if cmp, present := o.GetAttr("__gt__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return false //l  > r
}

// Lookup the equal operator and execute it, if one exists.
func (o *Object) Eq(l, r *Object) (bool) {
    if cmp, present := o.GetAttr("__eq__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return false //l  == r
}

// Lookup the not equal operator and execute it, if one exists.
func (o *Object) Neq(l, r *Object) (bool) {
    if cmp, present := o.GetAttr("__neq__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return false //l != r
}

// Lookup the less than or equal operator and execute it, if one exists.
func (o *Object) Lte(l, r *Object) (bool) {
    if cmp, present := o.GetAttr("__lte__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return false //l  <= r
}

// Lookup the greater than or equal operator and execute it, if one exists.
func (o *Object) Gte(l, r *Object) (bool) {
    if cmp, present := o.GetAttr("__gte__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return false //l  >= r
}

*/
