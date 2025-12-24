package expander

import "fmt"

// processingStack manages files currently being processed to detect circular references
type processingStack struct {
	files map[string]bool
}

// newProcessingStack creates a new processingStack
func newProcessingStack() *processingStack {
	return &processingStack{
		files: make(map[string]bool),
	}
}

// enter adds a file to the processing stack
// Returns error if the file is already being processed (circular reference)
func (ps *processingStack) enter(path string) error {
	if ps == nil {
		return fmt.Errorf("processing stack is nil")
	}
	if ps.files[path] {
		return fmt.Errorf("circular reference detected: %s", path)
	}
	ps.files[path] = true
	return nil
}

// exit removes a file from the processing stack
func (ps *processingStack) exit(path string) {
	if ps == nil || ps.files == nil {
		return // TODO: 異常終了するべき
	}
	delete(ps.files, path)
}
