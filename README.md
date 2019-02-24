gotestlint is a linter for Go source code
it shows functions not had been covered by tests

## Usage

invoke gotestlint with one argument, that's path to package
```gotestlint``` doesn't go through your folders (don't recursive)
one checks only one package

```
$ gotestlint .
lint.go:20:there're not have any tests for LintFiles
lint.go:47:there're not have any tests for String
```

if there're all functions covered by tests,
it won't show any messages

```
$ gotestlint $GOPATH/src/github.com/your_the_best_project
```
