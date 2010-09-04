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

package types

type PyObject struct {
    Name string
    Attrs map[string] *PyObject 
}

// Object attribute getting interface.
type Getter interface {
    GetAttr(name String) (*PyObject, bool)     
}

// Object attribute setting interface.
type Setter interface {
    SetAttr(name String, value *PyObject)     
}

// Object rich comparison interface
type RichComparer interface {
    Lt(l, r *PyObject) (bool)
    Gt(l, r *PyObject) (bool)
    Eq(l, r *PyObject) (bool)
    Neq(l, r *PyObject) (bool)
    Lte(l, r *PyObject) (bool)
    Gte(l, r *PyObject) (bool)
}

// Get the value of an object's attribute.
func (o *PyObject) GetAttr(name String) (value *PyObject, present bool) {
    value, present := o.Attrs[name]  
}

// Set the value of an object's attribute.
func (o *PyObject) SetAttr(name String, value *PyObject) {
    o.Attrs[name] = value  
}

// Lookup the less than operator and execute it, if one exists.
func (o *PyObject) Lt(l, r *PyObject) (bool) {
    if cmp, present := o.GetAttr("__lt__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return l  < r
}

// Lookup the greater than operator and execute it, if one exists.
func (o *PyObject) Gt(l, r *PyObject) (bool) {
    if cmp, present := o.GetAttr("__gt__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return l  > r
}

// Lookup the equal operator and execute it, if one exists.
func (o *PyObject) Eq(l, r *PyObject) (bool) {
    if cmp, present := o.GetAttr("__eq__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return l  == r
}

// Lookup the not equal operator and execute it, if one exists.
func (o *PyObject) Neq(l, r *PyObject) (bool) {
    if cmp, present := o.GetAttr("__neq__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return l != r
}

// Lookup the less than or equal operator and execute it, if one exists.
func (o *PyObject) Lte(l, r *PyObject) (bool) {
    if cmp, present := o.GetAttr("__lte__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return l  <= r
}

// Lookup the greater than or equal operator and execute it, if one exists.
func (o *PyObject) Gte(l, r *PyObject) (bool) {
    if cmp, present := o.GetAttr("__gte__"); present {
        // Execute the rich comparison operator.
        return false
    }
    
    // Default to comparing the pointer values.  This
    // is probably wrong;
    return l  >= r
}


