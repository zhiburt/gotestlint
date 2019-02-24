gotestlint is a linter for Go source code

it shows functions not had been covered by tests

## Usage

invoke gotestlint with one argument, that's path to package ```gotestlint``` doesn't go through your folders (don't recursive). The one checks only one package

```
$ gotestlint .
lint.go:20:function LintFiles is not covered any tests
lint.go:47:function String is not covered any tests
```

if there're all functions covered by tests,
it won't show any messages

```
$ gotestlint $GOPATH/src/github.com/your_the_best_project
```

you can add option ```nolint: gotestlint``` in comment your function
if function is marked such way, ```gotestlint``` ignores that function
