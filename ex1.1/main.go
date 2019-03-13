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
