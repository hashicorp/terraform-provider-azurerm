package validate

func CacheNFSExport(i interface{}, k string) (warnings []string, errs []error) {
	return absolutePath(i, k)
}
