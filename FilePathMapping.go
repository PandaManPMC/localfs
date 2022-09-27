package localfs

func init() {
	dirMap = make(map[string]*dirMapping)
}

type dirMapping struct {
	dirName string
	dir     string
	relDir  string
	count   uint32
}

var dirMap map[string]*dirMapping

func putDirMap(typePath string, dm *dirMapping) {
	dirMap[typePath] = dm
}

func getDirMap(typePath string) *dirMapping {
	return dirMap[typePath]
}
