[![Go Report Card](https://goreportcard.com/badge/github.com/zhiburt/gotestlint)](https://goreportcard.com/report/github.com/zhiburt/gotestlint)
[![GoDoc](https://godoc.org/github.com/zhiburt/gotestlint?status.svg)](https://godoc.org/github.com/zhiburt/gotestlint)


gotestlint is a linter for Go source code

it shows functions not had been covered by tests

## Usage

invoke gotestlint with one argument, that's path to package ```gotestlint``` doesn't go through your folders (don't recursive). The one checks only one package

```
$ gotestlint .
lint.go:20:function LintFiles is not covered any tests
lint.go:47:function String is not covered any tests
```

there's a recursive variant for it

```
gotestlint ./...
lint.go:27:function LintFiles is not covered any tests
```

It also can check only a bit files, something like this

```
gotestlint lint.go test_folder/test_folder.go 
lint.go:20:function LintSource is not covered any tests
lint.go:27:function LintFiles is not covered any tests
lint.go:54:function String is not covered any tests
test_folder/test_folder.go:3:function SomeFunc is not covered any tests
...
```

if all functions are covered by tests,
it won't show any messages

```
$ gotestlint $GOPATH/src/github.com/your_the_best_project
```

you can add option ```nolint: gotestlint``` in comment your function
if function is marked such way, ```gotestlint``` ignores that function,
for example

```
// Foo ...
// 
// nolint: gotestlint
func Foo(){
}
```
