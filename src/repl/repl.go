package repl

import (
	"Servo/src/lexer"
	"Servo/src/parser"
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

		io.WriteString(out, program.String())
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