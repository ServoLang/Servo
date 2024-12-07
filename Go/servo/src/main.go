package main

import (
	"Go/servo/src/lexer"
	"Go/servo/src/parser"
	"github.com/sanity-io/litter"
	"os"
)

func main() {
	bytes, _ := os.ReadFile("./examples/04.svo")
	tokens := lexer.Tokenize(string(bytes))

	ast := parser.Parse(tokens)
	litter.Dump(ast)
}
