package parser_test

import ( 
        "core/parser";
	    "testing";
	    "utf8"
)	

var test_string = "import err"    

func TestPeek(t *testing.T) {
	s, err := parser.Open("test_data/test1.py")
	
	if err!=nil {		
		t.Errorf("Open stream: %#v\n", err)
	}
	
	b, err := s.Peek()
	if err!=nil {
		t.Errorf("Peek() error: %#v", err)
	}
	
	if tb, _ := utf8.DecodeRuneInString(test_string); b!=tb {		
		t.Error("Expected to Peek() a(n) %#v", tb)
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

