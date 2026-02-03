package ancestor

import "maps"

// List of ancestor file paths referencing this file.
// for preventing circular references
type AncestorSet struct {
	fileAbsPaths map[string]struct{}
}

func NewEmptyAncestorSet() AncestorSet {
	return AncestorSet{
		fileAbsPaths: make(map[string]struct{}),
	}
}

// The presence of absPath in the ancestor set indicates a circular reference.
func (as AncestorSet) Contain(absPath string) bool {
	_, exists := as.fileAbsPaths[absPath]
	return exists
}

// Returns a new ancestorSet with absPath added. The original set remains unchanged.
func (as AncestorSet) Add(absPath string) AncestorSet {
	newPaths := make(map[string]struct{}, len(as.fileAbsPaths)+1)

	maps.Copy(newPaths, as.fileAbsPaths)

	newPaths[absPath] = struct{}{}

	return AncestorSet{
		fileAbsPaths: newPaths,
	}
}
