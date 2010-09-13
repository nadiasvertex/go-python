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

   The parser package implements a simple library for parsing EBNF
   grammars.
   
   The ast objects are the internal representation of the abstract syntax tree
   of the Python language.  These may be quite different than the CPython ast.
*/

package parser

type Ast interface {
    Next() Node*
    Prev() Node*
}

type Node struct {
    Parent  Ast*
    Op      int
}

type LiteralIntNode {
    *Node
    Value int
} 

type LiteralStringNode {
    *Node
    Value string
}



