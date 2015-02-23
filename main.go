package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/easonliao/gofp/ast"
	"github.com/easonliao/gofp/parser"
)

func main() {
	file := os.Stdin
	res := ast.NilObj

	if len(os.Args) > 1 {
		var err error
		file, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	reader := bufio.NewReader(file)
	sc := ast.NewScope(nil)

	for {
		if file == os.Stdin {
			fmt.Print(">")
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println(res)
				return
			}
			fmt.Println(err)
			return
		}
		expr, err := parser.ParseExpr([]byte(line))
		if err != nil {
			fmt.Println(err)
			continue
		}
		ast.Print(expr)
		res, err = expr.Eval(sc)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if file == os.Stdin {
			fmt.Println(res)
		}
	}
}
