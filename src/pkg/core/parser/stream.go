/* 
   The parser package implements a simple library for parsing EBNF
   grammars.
   
   The scanner, lexer, and parser are all implemented together for
   efficiency.  The scanner provides a useful stream paradigm that
   allows merging data on the fly into the stream.   
   
*/
package parser

import (
    "container/vector";
	"io/ioutil";
	"os";
	"strings";
)

// Provides stream context data for the scanner
type Stream struct {
	data  []int
	index int
	at    int
	row   uint
	col   uint
	
	streams *vector.Vector
	
	in_context_mgr   bool
	ctx_mgr_rollback bool	
}

// Initializes a new stream context object.
func newStream(data string) *Stream {
	return &Stream{strings.Runes(data), 0, -1, 1, 1, new(vector.Vector) , false, false}
}

// Open creates a stream by loading data from a file.  If the 
// os returns an error, the stream will be nil and the
// error will be returned in err.
func Open(filename string) (s *Stream, err os.Error) {

	contents, err := ioutil.ReadFile(filename)
	if err!= nil {
		return nil, err		
	}	
	
	return newStream(string(contents)), nil
}

func (s *Stream) Merge(data string, filename string) os.Error {

	return nil
}

func (s *Stream) Peek() (ch int, err os.Error) {

	at := s.at + 1
	data := s.data
		
	if at >= len(s.data) {
		index := s.index+1
		at = 0
		
		if index >= len(*s.streams) {
			ch   = 0
			err  = os.EOF
			return 
		}			 

		data = s.streams.At(index).(Stream).data
	}
	
	ch   = data[at]
	err  = nil
	 
	return
}

func (s *Stream) Read() (ch int, err os.Error) {

	s.at+=1
	
	if s.at>=len(s.data) {
		s.index+=1
		s.at = 0
		
		if s.index > len(*s.streams) {
			ch   = 0
			err  = os.EOF
			return 
		}
		
		s.data = s.streams.At(s.index).(Stream).data
	}
	
	if s.data[s.at]=='\n' {
		s.row+=1
		s.col=1
	} else {
		s.col+=1
	}
	
	ch   = s.data[s.at]
	err  = nil
	
	return
}

func (s *Stream) GetLoc() {
}

func (s *Stream) SetLoc() {
}

func (s *Stream) GetMarker() {
}

func (s *Stream) SetMarker() {
}

func (s *Stream) BeginTransaction() {
}

func (s *Stream) Commit() {
}

func (s *Stream) Rollback()	{
}
