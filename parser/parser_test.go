package parser_test

import (
	"path/filepath"
	"reflect"
	"testing"

	tf "github.com/gotestlint/iternal/testfolder"
	"github.com/gotestlint/parser"
)

func TestParse(t *testing.T) {
	const testfoldername = "testfolder"
	cases := []struct {
		path     string
		files    []*tf.File
		expected map[string][]byte
		err      error
	}{
		{
			path: ".",
			files: []*tf.File{
				&tf.File{Name: testfoldername, Content: []byte(""), Folder: []*tf.File{
					&tf.File{Name: "invalidfile.go", Content: []byte("package testfolder\n\nfunc Foo(){}"), Folder: nil},
					&tf.File{Name: "valid.go", Content: []byte("package testfolder\n\nfunc Foo2(){}"), Folder: nil},
					&tf.File{Name: "valid_test.go", Content: []byte("package testfolder\nimport \"testing\"\nfunc TestFoo2(t *testing.T){}"), Folder: nil},
				}},
			},
			expected: map[string][]byte{
				"valid.go":       []byte("package testfolder\n\nfunc Foo2(){}"),
				"valid_test.go":  []byte("package testfolder\nimport \"testing\"\nfunc TestFoo2(t *testing.T){}"),
				"invalidfile.go": []byte("package testfolder\n\nfunc Foo(){}"),
			},
			err: nil,
		},
		{
			path: ".",
			files: []*tf.File{
				&tf.File{Name: testfoldername, Content: []byte(""), Folder: []*tf.File{
					&tf.File{Name: "invalidfile", Content: []byte("package testfolder\n\nfunc Foo(){}"), Folder: nil},
				}},
			},
			expected: map[string][]byte{},
			err:      nil,
		},
		{
			path: ".",
			files: []*tf.File{
				&tf.File{Name: testfoldername, Content: []byte(""), Folder: []*tf.File{
					&tf.File{Name: "invalidfile.go", Content: []byte(""), Folder: nil},
				}},
			},
			expected: map[string][]byte{
				"invalidfile.go": []byte{},
			},
			err: nil,
		},
	}
	testfolder := tf.NewTestFolder()

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			_, err := testfolder.Create(c.path, c.files...)
			if err != nil {
				t.Fatal("Prepare folder error")
			}
			defer testfolder.Close()
			res, err := parser.Parse(filepath.Join(c.path, c.files[0].Name))
			if err != c.err {
				t.Errorf("[error] expexted error %v\nbut was %v", c.err, err)
			}
			if !reflect.DeepEqual(c.expected, res) {
				t.Errorf("[error] expexted map %#v\nbut was %#v", c.expected, res)
			}
		})
	}
}
