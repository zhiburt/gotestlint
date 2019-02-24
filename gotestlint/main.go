package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	lint "github.com/gotestlint"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("don't have enought arguments")
	}
	path := os.Args[1]

	files, err := parseFolder(path)
	_handleError(err)

	l := lint.Linter{}
	resp, err := l.LintFiles(files)
	_handleError(err)

	for _, problem := range resp {
		fmt.Println(problem)
	}
}

var ParseFolderErr = errors.New("can't parse folder")

func parseFolder(path string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	information, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, ParseFolderErr
	}

	for _, info := range information {
		if info.IsDir() {
			continue
		}
		if strings.HasSuffix(info.Name(), ".go") == false {
			continue
		}

		p := filepath.Join(path, info.Name())
		body, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, ParseFolderErr
		}
		files[info.Name()] = body
	}

	return files, err
}

func _handleError(err error) {
	if err != nil {
		log.Fatalf("erorr hapend %v", err)
	}
}
