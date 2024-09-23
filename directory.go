package main

import (
	"fmt"
	"sort"
	"strings"
)

const (
	SEPARATOR string = "/"
	INDENT    string = "  "
)

func getParentPath(path []string) []string {
	if len(path) == 0 {
		return path
	}

	return path[0 : len(path)-1]
}

func (root *Directory) get(path []string) *Directory {
	// the path is already validated, so traverse it
	dir := root
	ok := false
	for _, name := range path {
		if dir, ok = dir.children[name]; !ok {
			return nil
		}
	}

	return dir
}

type Directory struct {
	children map[string]*Directory
}

// NewDirectory creates a new, empty Directory
func NewDirectory() *Directory {
	return &Directory{children: map[string]*Directory{}}
}

// list all directories indented by hierarchy
func (root *Directory) List() { root.list(0) }
func (root *Directory) list(level int) {
	// sort the
	names := make([]string, 0, len(root.children))
	for n := range root.children {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
		dir := root.children[name]
		fmt.Printf("%s%s\n", strings.Repeat(INDENT, level), name)
		dir.list(level + 1)
	}
}

func (root *Directory) Create(path []string) error {
	errPrefix := fmt.Sprintf("Cannot create %s", strings.Join(path, SEPARATOR))

	// validate that the parent path exists and the target does not
	if err := NewValidator(false, TargetDoesNotExist).Validate(root, path); err != nil {
		return fmt.Errorf("%s - %s", errPrefix, err.Error())
	}

	parentDir := root.get(getParentPath(path))
	name := path[len(path)-1]
	parentDir.children[name] = NewDirectory()

	return nil
}

func (root *Directory) Delete(path []string) error {
	errPrefix := fmt.Sprintf("Cannot delete %s", strings.Join(path, SEPARATOR))

	if err := NewValidator(false, TargetDoesExist).Validate(root, path); err != nil {
		return fmt.Errorf("%s - %s", errPrefix, err.Error())
	}

	parentDir := root.get(getParentPath(path))
	name := path[len(path)-1]
	delete(parentDir.children, name)

	return nil
}

func (root *Directory) Move(src []string, dst []string) error {
	errPrefix := fmt.Sprintf(
		"Cannot move %s to %s",
		strings.Join(src, SEPARATOR),
		strings.Join(dst, SEPARATOR))

	// src must exist
	if err := NewValidator(false, TargetDoesExist).Validate(root, src); err != nil {
		return fmt.Errorf("%s - %s", errPrefix, err.Error())
	}

	// dst must exist, but could be the empty path (/)
	if err := NewValidator(true, TargetDoesExist).Validate(root, dst); err != nil {
		return fmt.Errorf("%s - %s", errPrefix, err.Error())
	}

	srcDir := root.get(src)
	dstDir := root.get(dst)
	srcParent := root.get(getParentPath(src))
	srcName := src[len(src)-1]

	delete(srcParent.children, srcName)
	dstDir.children[srcName] = srcDir

	return nil
}
