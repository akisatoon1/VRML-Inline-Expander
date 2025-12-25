package expander

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
	newFiles := make(map[string]struct{}, len(as.fileAbsPaths)+1)

	for k, v := range as.fileAbsPaths {
		newFiles[k] = v
	}

	newFiles[absPath] = struct{}{}

	return ancestorSet{
		fileAbsPaths: newFiles,
	}
}
