package validate

func LogAnalyticsLinkedStorageAccountWorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	return logAnalyticsGenericName(i, k)
}
