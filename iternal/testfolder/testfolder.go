package testfolder

import (
	"log"
	"os"
	"path/filepath"
)

type (
	TestFolder interface {
		Create(string, ...*File) ([]*File, error)
		Close() error
	}

	File struct {
		Name    string
		Content []byte
		Folder  []*File
	}

	DefaultTestFolder struct {
		roots []*File
		path  string
	}
)

func NewTestFolder() TestFolder {
	return &DefaultTestFolder{}
}

func (f *File) IsFolder() bool {
	return f.Folder != nil
}

func (folder *DefaultTestFolder) Create(path string, files ...*File) ([]*File, error) {
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

func (file *DefaultTestFolder) Close() error {
	for _, root := range file.roots {
		path := filepath.Join(file.path, root.Name)
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
