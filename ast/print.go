package ast

import (
	"fmt"
	"reflect"

	"github.com/easonliao/gofp/token"
)

func Print(root Expr) {
	var p printer
	p.print(reflect.ValueOf(root))
}

type printer struct {
	indent int
}

func (p *printer) print(x reflect.Value) {
	//fmt.Println("[", x.Kind(), "]")
	switch x.Kind() {
	case reflect.Interface:
		p.print(x.Elem())
	case reflect.Ptr:
		p.print(x.Elem())
	case reflect.Struct:
		t := x.Type()
		p.printf("%s {", t)
		p.printf("\n")
		p.indent++
		for i, n := 0, t.NumField(); i < n; i++ {
			// exclude non-exported fields because their
			// values cannot be accessed via reflection
			name := t.Field(i).Name
			value := x.Field(i)
			p.printfWithIndent("%s: ", name)
			p.print(value)
		}
		p.indent--
		p.printfWithIndent("}")
		p.printf("\n")
	case reflect.Slice:
		p.printf("%s (len = %d) {", x.Type(), x.Len())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printfWithIndent("%d: ", i)
				p.print(x.Index(i))
			}
			p.indent--
		}
		p.printfWithIndent("}")
		p.printf("\n")
	default:
		v := x.Interface()
		switch v := v.(type) {
		case string:
			p.printf("%q", v)
			p.printf("\n")
			return
		case float64, float32:
			p.printf("%f", v)
			p.printf("\n")
			return
		case token.Token:
			p.printf("%q", token.TokenName(v))
			p.printf("\n")
			return
		}
	}
}

func (p *printer) printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (p *printer) printfWithIndent(format string, args ...interface{}) {
	for i := 0; i < p.indent; i++ {
		fmt.Print(".  ")
	}
	p.printf(format, args...)
}
