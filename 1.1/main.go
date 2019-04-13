/*
Modify the echo program to also print os.Args[0], the name of the command that invoked it.
*/
package main

import "os"

func main() {
	s, sep := "", ""
	for _, v := range os.Args {
		s += sep + v
		sep = " "
	}
	println(s)
}
