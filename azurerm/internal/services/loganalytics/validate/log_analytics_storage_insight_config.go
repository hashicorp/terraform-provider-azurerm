package validate

func LogAnalyticsStorageInsightConfigName(i interface{}, k string) (warnings []string, errors []error) {
	return LogAnalyticsGenericName(i, k)
}

func LogAnalyticsStorageInsightConfigWorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	return LogAnalyticsGenericName(i, k)
}
