package testfolder

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCreate(t *testing.T) {
	const testfoldername = "test"
	cases := []struct {
		path     string
		files    []*File
		expected int
		err      error
	}{
		{
			path: ".",
			files: []*File{
				&File{Name: testfoldername, Content: []byte(""), Folder: []*File{
					&File{Name: "invalidfile.go", Content: []byte("package testfolder\n\nfunc Foo(){}"), Folder: nil},
					&File{Name: "valid.go", Content: []byte("package testfolder\n\nfunc Foo2(){}"), Folder: nil},
					&File{Name: "valid_test.go", Content: []byte("package testfolder\nimport \"testing\"\nfunc TestFoo2(t *testing.T){}"), Folder: nil},
				}},
			},
			expected: 3,
			err:      nil,
		},
		{
			path: ".",
			files: []*File{
				&File{Name: testfoldername, Content: []byte(""), Folder: []*File{}},
			},
			expected: 0,
			err:      nil,
		},
	}

	testfolder := NewTestFolder()
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			_, err := testfolder.Create(c.path, c.files...)
			if err != nil {
				t.Fatal("Prepare folder error")
			}
			defer testfolder.Close()

			if amaunt, err := checkfolder(filepath.Join(c.path, c.files[0].Name)); err != c.err {
				t.Errorf("[error] expexted error %v\nbut was %v", c.err, err)
			} else if amaunt != c.expected {
				t.Errorf("[error] expexted amaunt %v\nbut was %v", c.expected, amaunt)
			}
		})
	}
}

func checkfolder(path string) (int, error) {
	info, err := ioutil.ReadDir(path)
	if err == nil {
		var amaunt = 0
		for _, i := range info {
			n, err := checkfolder(filepath.Join(path, i.Name()))
			if err != nil {
				return n, err
			}
			amaunt += n
		}

		return amaunt, err
	}
	_, err = ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	return 1, nil
}
