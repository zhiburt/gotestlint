package lint

import (
	"reflect"
	"testing"
)

func TestLintSource(t *testing.T) {
	cases := []struct {
		files    map[string][]byte
		expected []string
	}{
		{
			files: map[string][]byte{
				"t.go":      []byte("package t\n\nfunc FirstFoo(){\n}\n\nfunc SecondFoo(){\n}"),
				"t_test.go": []byte("package t\n\nfunc TestFirstFoo(t *testing.T){\n}\n\nfunc TestSecondFoo(t *testing.T){\n}"),
			},
			expected: nil,
		},
		{
			files: map[string][]byte{
				"t.go": []byte("package t\n\nfunc FirstFoo(){\n}\n\nfunc SecondFoo(){\n}"),
			},
			expected: []string{
				"t.go:3:there're not have any tests for FirstFoo",
				"t.go:6:there're not have any tests for SecondFoo",
			},
		},
		{
			files: map[string][]byte{
				"t.go":      []byte("package t\n\nfunc FirstFoo(){\n}\n\nfunc SecondFoo(){\n}"),
				"t_test.go": []byte("package t\n\nfunc TestSecondFoo(t *testing.T){\n}"),
			},
			expected: []string{
				"t.go:3:there're not have any tests for FirstFoo",
			},
		},
		{
			files: map[string][]byte{
				"t.go":      []byte("package t\n\n// FirstFoo\n// nolint: gotestlint\nfunc FirstFoo(){\n}\n\nfunc SecondFoo(){\n}"),
				"t_test.go": []byte("package t\n\nfunc TestSecondFoo(t *testing.T){\n}"),
			},
			expected: nil,
		},
	}

	linter := Linter{}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			act, err := linter.LintFiles(c.files)
			if err != nil {
				t.Error("unexpected erorr", err)
			}
			if !reflect.DeepEqual(c.expected, adviseToString(act)) {
				t.Errorf("mistake  expected: %v\nbut was: %v", c.expected, act)
			}
		})
	}
}

func adviseToString(advs []Advise) []string {
	var str []string
	for _, problem := range advs {
		str = append(str, problem.String())
	}

	return str
}
