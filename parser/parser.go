/*Package parser parses folders*/
package parser

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	// ErrParseFolder mean that erorr hapened when was reading a folder
	ErrParseFolder = errors.New("can't parse folder")

	// ErrNoGOFile mean that you try parse no golang file
	ErrNoGOFile = errors.New("can't parse no .go file")

	// ErrParseFile mean that erorr hapened when was reading a file
	ErrParseFile = errors.New("can't parse file")
)

// RECURSIVEPATH represents marker for recursive call on all go files
const RECURSIVEPATH = "./..."

// Parse function parse folder/file on path
// return map where key is name of file, value is content and error
func Parse(path string) (map[string][]byte, error) {
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
		return nil, ErrParseFolder
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
