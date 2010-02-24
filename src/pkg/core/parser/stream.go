/* 
   The parser package implements a simple library for parsing EBNF
   grammars.
   
   The scanner, lexer, and parser are all implemented together for
   efficiency.  The scanner provides a useful stream paradigm that
   allows merging data on the fly into the stream.   
   
*/
package parser

import (
    "container/list";
	"io/ioutil";
	"fmt";
	"os";
	"strings";
)

// Provides stream data for the scanner
type Stream struct {		
	// The list of context objects.
	streams *list.List
	
	// The current context object.
	cur     *list.Element;	
	
	
}

// Initializes a new stream object.
func newStream() *Stream {
	return &Stream{new(list.List), nil}
}

// Holds context data about individual streams
type context struct {
	data     []int
	at       int
	row      uint
	col      uint
	name     string	
}

// Initializes a new stream context object.
func newStreamContext(data string, name string) *context {
	return &context{strings.Runes(data), -1, 1, 1, name }
}

// Splits a stream context into two pieces.  The first
// piece is the context passed in at the current read
// point.  The second is a new context with all of the
// data after the current read point.  If the read point
// is at the end of the data, nil will be returned.
func splitStreamContext(ctx *context) *context {

	partition := ctx.at+1			
	if partition >= len(ctx.data) {
		return nil
	}
	
	// Slice the current data into two pieces
	left  := ctx.data[0:partition]
	right := ctx.data[partition:]
	
	// Update the current context.
	ctx.data = left
		
	// Return a new context.
	return &context{right, -1, ctx.row, ctx.col, ctx.name}
}


// Open creates a stream by loading data from a file.  If the 
// os returns an error, the stream will be nil and the
// error will be returned in err.
func Open(filename string) (s *Stream, err os.Error) {

	contents, err := ioutil.ReadFile(filename)
	if err!= nil {
		return nil, err		
	}	
	
	s = newStream()
	s.MergeFromString(string(contents), filename)	
			
	return 
}

// Merge provides a way to save the current stream state and switch
// to an alternate data input.  This is useful for situations when
// you need to process an "include"-style directive, or you want to
// insert alternate data into the stream.
func (s *Stream) MergeFromString(data string, name string) {

	// If the data to merge is empty, do nothing
	if len(data)==0 || s==nil {
		return
	}		
	
	
	switch s.cur {
	// If our current context is empty, insert
	// a new context onto the stack.  It will get
	// switched to on the next Read() or Peek()
	case nil:
		s.cur = s.streams.PushFront(newStreamContext(data, name))				
	
	default:	
		// Get the current context.
		ctx := s.cur.Value.(* context)

		switch {	
		// If the current stream is at the end, insert a new stream
		// after it.	
		case ctx.at >= len(ctx.data):
			s.streams.InsertAfter(newStreamContext(data, name), s.cur)

		// Split the current stream.  Insert the merge stream 
		// in the middle.	
		default:						
			merge_stream := newStreamContext(data, name)
	
			// Insert a new stream object with the merged data after the current stream.
			merge_stream_el := s.streams.InsertAfter(merge_stream, s.cur)
	
			// Insert the split stream data after the merged stream.
			if split_stream := splitStreamContext(ctx); split_stream!=nil && merge_stream_el!=nil {
				s.streams.InsertAfter(split_stream, merge_stream_el)
			}
	  	}		
	}   

}

// Peek gives you the next character (not byte) from the
// data stream.  If the data stream is empty it will 
// return a 0 as the character and set err to os.EOF.
// Consecutive exeutions of Peek will return the same
// value.
func (s *Stream) Peek() (ch int, err os.Error) {
	
	if s==nil || s.cur==nil {
		return 0, os.EOF
	}

	ctx := s.cur.Value.(*context)

	at := ctx.at + 1
	data := ctx.data
		
	if at >= len(ctx.data) {
		next := s.cur.Next()
		at = 0
		
		if next==nil {
			ch   = 0
			err  = os.EOF
			return 
		}			 

		data = s.cur.Value.(*context).data
	}
	
	ch   = data[at]
	err  = nil
	 
	return
}

// Read consumes the next character (not byte) from the
// data stream.  If the data stream is empty it will 
// return a 0 as the character and set err to os.EOF.
func (s *Stream) Read() (ch int, err os.Error) {

	if s==nil || s.cur==nil {
		return 0, os.EOF
	}

	ctx := s.cur.Value.(*context)

	ctx.at+=1
	
	if ctx.at>=len(ctx.data) {	
		next := s.cur.Next()		
		
		if next == nil {
			ch   = 0
			err  = os.EOF
			return 
		}
		
		s.cur = next
		ctx   = s.cur.Value.(*context)
		
		// Increment as we would if we had just entered the Read() block.
		ctx.at+=1
	}
	
	if ctx.data[ctx.at]=='\n' {
		ctx.row+=1
		ctx.col=1
	} else {
		ctx.col+=1
	}
	
	ch   = ctx.data[ctx.at]
	err  = nil
	
	return
}

// GetLoc returns the row and column of 
// the read head for the Stream.  If the
// Stream is empty, GetLoc will return 
// 0,0.
func (s *Stream) GetLoc() (uint, uint) {
	if s==nil || s.cur == nil {
		return 0, 0
	}
	
	ctx := s.cur.Value.(*context)
	
	return ctx.row, ctx.col
}

// SetLoc allows you to override the current
// row and column locations in the Stream.  This
// doesn't actually move the read head, it just
// changes where the reported values from GetLoc
// are.  
func (s *Stream) SetLoc(row, col uint) {
	if s==nil || s.cur == nil {
		return
	}
	
	ctx := s.cur.Value.(*context)
	
	ctx.row = row
	ctx.col = col
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

// DumpStreamContext will stringify a Stream object.
func (s *Stream) DumpStreamContext() string {
	return fmt.Sprintf("%T %+v", s, s)
}

