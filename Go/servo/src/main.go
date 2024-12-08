package main

import (
	"Go/servo/src/lexer"
	"Go/servo/src/parser"
	"github.com/sanity-io/litter"
	"os"
)

func main() {
	file := "./examples/05.svo"
	bytes, _ := os.ReadFile(file)
	tokens := lexer.Tokenize(string(bytes))

	ast := parser.Parse(tokens, file)
	litter.Dump(ast)
}
