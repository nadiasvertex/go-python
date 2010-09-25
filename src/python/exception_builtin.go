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

type BaseExceptionObject struct {
    ObjectData
    args Object 
}

func (e *BaseExceptionObject) GetAttr(name string) (value Object, present bool) {
    if "args"==name {
        return e.args, true    
    }
    
    return nil, false
}