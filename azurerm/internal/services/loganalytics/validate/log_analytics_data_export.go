package validate

func LogAnalyticsDataExportWorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	return LogAnalyticsGenericName(i, k)
}

func LogAnalyticsDataExportName(i interface{}, k string) (warnings []string, errors []error) {
	return LogAnalyticsGenericName(i, k)
}
