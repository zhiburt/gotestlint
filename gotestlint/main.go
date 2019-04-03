package main

import (
	"fmt"
	"log"
	"os"

	lint "github.com/gotestlint"
	"github.com/gotestlint/parser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("don't have enought arguments")
	}

	for _, path := range os.Args[1:] {
		files, err := parser.Parse(path)
		_handleError(err)

		l := lint.Linter{}
		resp, err := l.LintFiles(files)
		_handleError(err)

		for _, problem := range resp {
			fmt.Println(problem)
		}
	}
}

func _handleError(err error) {
	if err != nil {
		log.Fatalf("erorr hapend %v", err)
	}
}
