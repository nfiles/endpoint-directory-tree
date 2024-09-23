package main

import (
	"errors"
	"fmt"
)

type TargetStatus int

const (
	TargetDoesNotExist TargetStatus = iota
	TargetOptional
	TargetDoesExist
)

type Validator struct {
	allowEmptyPath bool
	targetStatus   TargetStatus
}

func NewValidator(allowEmptyPath bool, targetStatus TargetStatus) Validator {
	return Validator{
		allowEmptyPath: allowEmptyPath,
		targetStatus:   targetStatus,
	}
}

func (v Validator) Validate(root *Directory, path []string) error {
	if !v.allowEmptyPath && len(path) == 0 {
		return errors.New("empty path")
	}

	dir := root
	exists := false

	for i, name := range path {
		if name == "" {
			return errors.New("empty directory name")
		}

		isTarget := i == len(path)-1
		dir, exists = dir.children[name]
		if !exists {
			// intermediate dirs are always required
			if !isTarget || v.targetStatus == TargetDoesExist {
				return fmt.Errorf("%s does not exist", name)
			}
		} else {
			if isTarget && v.targetStatus == TargetDoesNotExist {
				return fmt.Errorf("%s already exists", name)
			}
		}
	}

	return nil
}
