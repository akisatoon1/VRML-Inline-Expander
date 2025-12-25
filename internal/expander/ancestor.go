package expander

import "maps"

// List of ancestor file paths referencing this file.
// for preventing circular references
type ancestorSet struct {
	fileAbsPaths map[string]struct{}
}

func newEmptyAncestorSet() ancestorSet {
	return ancestorSet{
		fileAbsPaths: make(map[string]struct{}),
	}
}

// The presence of absPath in the ancestor set indicates a circular reference.
func (as ancestorSet) contains(absPath string) bool {
	_, exists := as.fileAbsPaths[absPath]
	return exists
}

// Returns a new ancestorSet with absPath added. The original set remains unchanged.
func (as ancestorSet) add(absPath string) ancestorSet {
	newPaths := make(map[string]struct{}, len(as.fileAbsPaths)+1)

	maps.Copy(newPaths, as.fileAbsPaths)

	newPaths[absPath] = struct{}{}

	return ancestorSet{
		fileAbsPaths: newPaths,
	}
}
