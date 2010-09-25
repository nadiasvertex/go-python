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

  Test for the static single assignment subsystem.   
  
*/

package python

import (        
        "testing"            
)


func TestWriteElements(t *testing.T) {    
    ctx := new (SsaContext)
    ctx.Init()
    
    el := new (SsaElement)
    
    for i:=0; i<256; i++ {
        if ctx.Write(el) != i {
            t.Errorf("Write returned an incorrect value as the new id.\n")  
        }
    }    
}
