package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func main() {
	// 源代码字符串
	src := `
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`

	// 设置文件解析的文件集
	fset := token.NewFileSet()
	// 解析源代码文件
	node, err := parser.ParseFile(fset, "example.go", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// 遍历 AST
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			fmt.Printf("Function: %s\n", x.Name)
		case *ast.CallExpr:
			if len(x.Args) > 0 {
				fmt.Printf("Call to Function: %s\n", x.Fun)
			}
		}

		// 继续遍历子节点
		return true
	})

	// 配置打印机
	config := printer.Config{
		Mode:   printer.UseSpaces | printer.TabIndent | printer.SourcePos,
		Indent: 4,
	}

	// 创建一个缓冲区以捕获输出
	var buf bytes.Buffer
	if err := config.Fprint(&buf, fset, node); err != nil {
		panic(err)
	}

	// 将 AST 输出到控制台
	bufio.NewReader(&buf).WriteTo(os.Stdout)
}
