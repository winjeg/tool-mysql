package main

import (
	"fmt"
	"github.com/winjeg/tool-mysql/parser"
)

func main() {
	astNode, err := parser.Parse("SELECT a, b FROM t")
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	}
	fmt.Printf("%v\n", parser.Extract(astNode))
}