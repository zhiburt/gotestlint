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

	for _, path := range os.Args[1:] {
		files, err := parseFolder(path)
		_handleError(err)

		l := lint.Linter{}
		resp, err := l.LintFiles(files)
		_handleError(err)

		for _, problem := range resp {
			fmt.Println(problem)
		}
	}
}

var ParseFolderErr = errors.New("can't parse folder")

const RECURSIVEPATH = "./..."

func parseFolder(path string) (map[string][]byte, error) {
	if path == RECURSIVEPATH {
		return parse(".", true)
	}

	return parse(path, false)
}

func parse(path string, recursive bool) (map[string][]byte, error) {
	files := make(map[string][]byte)
	information, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, ParseFolderErr
	}

	for _, info := range information {
		p := filepath.Join(path, info.Name())
		if info.IsDir() {
			if recursive {
				m, err := parse(p, recursive)
				if err != nil {
					return nil, err
				}

				copyMapToMap(files, m)
			}
			continue
		}
		if strings.HasSuffix(info.Name(), ".go") == false {
			continue
		}

		body, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, ParseFolderErr
		}
		files[info.Name()] = body
	}

	return files, err
}

func copyMapToMap(lhs, rhs map[string][]byte) {
	for key, value := range rhs {
		lhs[key] = value
	}
}

func _handleError(err error) {
	if err != nil {
		log.Fatalf("erorr hapend %v", err)
	}
}
