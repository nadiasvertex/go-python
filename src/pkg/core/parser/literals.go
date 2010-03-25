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
*/

package parser


//func Any_literal(acceptable string, s Stream, log Log) {
    /* Matches patterns like [a-zA-Z] or [0-9] */
  /*  
    lit  := ""
    cont := True
    loc  := s.GetLoc()
    
    with s:
       while cont:                
           c=s.peek()           
           if (c!=None) and (c in acceptable):
            lit+=s.read()
           else:
               if len(lit)==0:
                s.rollback()
                return None
               else:
                cont=False
           
    return { "value" : lit, "loc" : loc }
}*/
