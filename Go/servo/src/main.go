package main

import (
	"Go/servo/src/parser"
	"fmt"
	"github.com/sanity-io/litter"
	"os"
	"time"
)

func main() {
	file := "./examples/Test.svo"
	sourceBytes, _ := os.ReadFile(file)
	source := string(sourceBytes)
	start := time.Now()
	ast := parser.Parse(source)
	duration := time.Since(start)

	litter.Dump(ast)
	fmt.Printf("Duration: %v\n", duration)
}
