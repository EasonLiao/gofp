package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/easonliao/gofp/ast"
	"github.com/easonliao/gofp/parser"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	sc := ast.NewScope(nil)

	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		expr, err := parser.ParseExpr([]byte(line))
		if err != nil {
			fmt.Println(err)
			continue
		}

		obj, err := expr.Eval(sc)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(obj)
	}
}
