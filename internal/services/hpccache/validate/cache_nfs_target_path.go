package validate

func CacheNFSTargetPath(i interface{}, k string) (warnings []string, errs []error) {
	return relativePath(i, k)
}
