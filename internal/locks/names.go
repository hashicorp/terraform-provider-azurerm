package locks

func getLockName(name string, resourceType string) string {
	return resourceType + "." + name
}
