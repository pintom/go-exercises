package main

import (
	"fmt"
	"os"
)

func main() {
	s := ""
	for i, v := range os.Args {
		s += fmt.Sprintf("[%d]:%v\n", i, v)
	}
	println(s)
}
