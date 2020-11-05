package validate

func LogAnalyticsStorageInsightsName(i interface{}, k string) (warnings []string, errors []error) {
	return LogAnalyticsGenericName(i, k)
}

func LogAnalyticsStorageInsightsWorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	return LogAnalyticsGenericName(i, k)
}
