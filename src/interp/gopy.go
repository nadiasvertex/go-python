package main

import (
	"fmt"
	"flag"
)

var verbose_output = flag.Bool("v", false, "verbose output")
var show_version = flag.Bool("V", false, "show version information and exit")

func main() {
	flag.Parse()
	
	if *show_version {
		fmt.Printf("gopython version 0.1\n")
	}		
}
	
