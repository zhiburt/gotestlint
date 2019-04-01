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
		files, err := parse(path)
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
var ErrNoGOFile = errors.New("can't parse no .go file")
var ErrParseFile = errors.New("can't parse file")

const RECURSIVEPATH = "./..."

func parse(path string) (map[string][]byte, error) {
	if path == RECURSIVEPATH {
		return parseFolder(".", true)
	}

	if body, err := parseFile(path); err == nil {
		return map[string][]byte{path: body}, nil
	}

	return parseFolder(path, false)
}

func parseFolder(path string, recursive bool) (map[string][]byte, error) {
	files := make(map[string][]byte)
	information, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, ParseFolderErr
	}

	for _, info := range information {
		p := filepath.Join(path, info.Name())
		if info.IsDir() {
			if recursive {
				m, err := parseFolder(p, recursive)
				if err != nil {
					return nil, err
				}

				copyMapToMap(files, m)
			}
			continue
		}

		body, err := parseFile(p)
		if err != nil {
			if err == ErrNoGOFile {
				continue
			}

			return nil, err
		}
		files[info.Name()] = body
	}

	return files, err
}

func parseFile(path string) ([]byte, error) {
	if strings.HasSuffix(path, ".go") == false {
		return nil, ErrNoGOFile
	}

	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, ErrParseFile
	}

	return body, nil
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
