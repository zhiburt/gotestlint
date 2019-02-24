package lint

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

type Linter struct {
}

func (l *Linter) LintSource(filename string, source []byte) ([]Advise, error) {
	return l.LintFiles(map[string][]byte{filename: source})
}

func (l *Linter) LintFiles(files map[string][]byte) ([]Advise, error) {
	pkg := pkg{
		fset:  token.NewFileSet(),
		files: make(map[string]*file),
	}
	for name, src := range files {
		f, err := parser.ParseFile(pkg.fset, name, src, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		pkg.files[name] = &file{f: f, filename: name}
	}

	return pkg.lint()
}

type Advise struct {
	fName    string
	position token.Position
	file     *file
}

func (a Advise) String() string {
	return fmt.Sprintf("%s:%d:function %s is not covered any tests", a.file.filename, a.position.Line, a.fName)
}

type pkg struct {
	fset  *token.FileSet
	files map[string]*file
}

func (p *pkg) lint() ([]Advise, error) {
	var advs []Advise
	for _, file := range p.files {
		for _, f := range file.exportedFuncs() {
			if f.isTestFunc() {
				continue
			}
			if f.isNolint() {
				continue
			}

			if tstfile := p.testFileFor(f.file); tstfile != nil {
				if _, err := tstfile.funcWithPrefix("Test" + f.f.Name.Name); err != nil {
					advs = append(advs, Advise{
						fName:    f.f.Name.Name,
						file:     f.file,
						position: p.fset.Position(f.pos),
					})
				}
			} else {
				advs = append(advs, Advise{
					fName:    f.f.Name.Name,
					file:     f.file,
					position: p.fset.Position(f.pos),
				})
			}
		}
	}

	return advs, nil
}

func (p *pkg) testFileFor(f *file) *file {
	for filename := range p.files {
		if strings.HasPrefix(filename, f.f.Name.Name) && strings.HasSuffix(filename, "_test.go") {
			return p.files[filename]
		}
	}

	return nil
}

type file struct {
	f        *ast.File
	filename string // with sufix
}

func (f *file) funcWithPrefix(prefix string) (exportFunc, error) {
	for _, foo := range f.exportedFuncs() {
		if strings.HasPrefix(foo.f.Name.Name, prefix) {
			return foo, nil
		}
	}

	return exportFunc{}, errors.New("not foud this function in exports ones")
}

func (f *file) findExportFunc(fname string) (exportFunc, error) {
	for _, foo := range f.exportedFuncs() {
		if foo.f.Name.Name == fname {
			return foo, nil
		}
	}

	return exportFunc{}, errors.New("not foud this function in exports ones")
}

func (f *file) exportedFuncs() []exportFunc {
	var funcs []exportFunc
	for _, dec := range f.f.Decls {
		if fun, ok := dec.(*ast.FuncDecl); ok {
			if fun.Name.IsExported() {
				funcs = append(funcs, exportFunc{
					file: f,
					f:    fun,
					pos:  fun.Pos(),
				})
			}
		}
	}

	return funcs
}

type exportFunc struct {
	f    *ast.FuncDecl
	file *file
	pos  token.Pos
}

func (foo *exportFunc) isTestFunc() bool {
	return strings.HasPrefix(foo.f.Name.Name, "Test")
}

func (foo *exportFunc) isNolint() bool {
	if foo.f.Doc == nil {
		return false
	}

	r := regexp.MustCompile(`\s+nolint:.*gotestlint`)
	for _, comment := range foo.f.Doc.List {
		if r.MatchString(comment.Text) {
			return true
		}
	}

	return false
}
