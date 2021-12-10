package validate

func CacheNamespacePath(i interface{}, k string) (warnings []string, errs []error) {
	return absolutePath(i, k)
}
