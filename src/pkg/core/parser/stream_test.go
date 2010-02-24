package parser_test

import ( 
        "core/parser";
	    "testing";
	    "utf8"
)	

var test_string = "import err"    
var test_merge_data = "this and that"

func TestPeek(t *testing.T) {
	s, err := parser.Open("test_data/test1.py")
		
	if err!=nil {		
		t.Errorf("Open stream: %+v\n", err)
	}
		
	
	b, err := s.Peek()
	if err!=nil {
		t.Logf("Stream: %s", s.DumpStreamContext());
		t.Errorf("Peek() error: %+v", err)
	}
	
	if tb, _ := utf8.DecodeRuneInString(test_string); b!=tb {		
		t.Errorf("Expected to Peek() a(n) %#v", tb)
	}		
}

func TestRead(t *testing.T) {
		
	s, err := parser.Open("test_data/test1.py")
	
	if err!=nil {		
		t.Errorf("Open stream: %#v\n", err)
	}	
	
	for pos, tb := range test_string {
		b, err := s.Read()
		
		if err!=nil {
			t.Error("Read() error: %#v at index %d", err, pos)
		}
	
		if b != tb {		
			t.Errorf("Expected to Read() a(n) %#v but read a(n) %#v at index %d", tb, b, pos)
		}
	}
}

func TestMergeWithSplit(t *testing.T) {
		
	s, err := parser.Open("test_data/test1.py")
	
	if err!=nil {		
		t.Errorf("Open stream: %#v\n", err)
	}	
	
	// Read some data.
	for i := 0; i<5; i++ {
		s.Read()	
	}
	
	// Merge new data
	s.MergeFromString(test_merge_data, "my_test_data")
	
	t.Log("Trying to read merged data.")
	
	// Test that we can read the merged data.
	for pos, tb := range test_merge_data {
		b, err := s.Read()
		
		if err!=nil {
			t.Error("Read() error: %#v at index %d", err, pos)
		}
	
		if b != tb {		
			t.Errorf("Expected to Read() a(n) %#v but read a(n) %#v at index %d", tb, b, pos)
		}
	}
	
	t.Log("Trying to read more of the previously tested data.")
	
	// Test that we drop back to the previous data.
	for pos, tb := range test_string[5:] {
		b, err := s.Read()
		
		if err!=nil {
			t.Error("Read() error: %#v at index %d", err, pos)
		}
	
		if b != tb {		
			t.Errorf("Expected to Read() a(n) %#v but read a(n) %#v at index %d", tb, b, pos)
		}
	}
	
}


