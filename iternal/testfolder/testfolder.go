package testfolder

import (
	"log"
	"os"
	"path/filepath"
)

type (
	// TestFolder is represent of test folder which can be deleted with all its files
	// method Create() create all files if it's possible
	// method Close() delete all created files
	TestFolder interface {
		Create(string, ...*File) ([]*File, error)
		Close() error
	}

	// File is file of directory
	// if Folder == nil this is file, else folder
	File struct {
		Name    string
		Content []byte
		Folder  []*File
	}

	// defaultTestFolder implementation of TestFolder
	defaultTestFolder struct {
		roots []*File
		path  string
	}
)

// NewTestFolder returns TestFolder obj
func NewTestFolder() TestFolder {
	return &defaultTestFolder{}
}

// IsFolder checks this is folder or not
func (f *File) IsFolder() bool {
	return f.Folder != nil
}

// Create creates all files in path
// that's does it recursively
func (folder *defaultTestFolder) Create(path string, files ...*File) ([]*File, error) {
	folder.path = path
	folder.roots = files
	var created []*File
	for _, file := range files {
		createdfile, err := createfile(path, file)
		if err != nil {
			log.Printf("func [DefaultTestFolder.Create] error hapend %v", err)
			return nil, err
		}

		created = append(created, createdfile)
	}

	return created, nil
}

// Close removes all had created files
func (folder *defaultTestFolder) Close() error {
	for _, root := range folder.roots {
		path := filepath.Join(folder.path, root.Name)
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func createfile(path string, f *File) (*File, error) {
	pathto := filepath.Join(path, f.Name)
	if f.IsFolder() {
		err := os.MkdirAll(pathto, os.ModePerm)
		if err != nil {
			return nil, err
		}
		for _, file := range f.Folder {
			_, err := createfile(pathto, file)
			if err != nil {
				return nil, err
			}
		}

		return f, nil
	}

	file, err := os.Create(pathto)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Write(f.Content)
	if err != nil {
		return nil, err
	}

	return f, nil
}
