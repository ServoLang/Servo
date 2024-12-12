package repl

import (
	"Servo/src/compiler"
	"Servo/src/lexer"
	"Servo/src/object"
	"Servo/src/parser"
	"Servo/src/vm"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

const GEAR = `
                             ..                             
                         ..........                         
                .        ..........        .                
             ......   ................   ......             
            ....................................            
           ......................................           
             ......                      .......            
             ....                          ....             
            ....                            ....            
      .........       :.            .:       .........      
      .........       :.            .:       .........      
     .........        :.            ..        ........      
      .........                              .........      
      ...........                          ...........      
            ....................................            
             ..................................             
             ...................................            
           ......................................           
           ......................................           
             ......   ................   ......             
                .        ..........        .                
                         ..........                         
                            ....`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Compilation failed:\n %s\n", err)
			continue
		}

		code := comp.ByteCode()
		constants = code.Constants

		machine := vm.NewWithGlobalStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Machine failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect()) // TODO: Remove to force printing
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, GEAR)
	io.WriteString(out, "\nWhoops! Errors really grind my gears!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
