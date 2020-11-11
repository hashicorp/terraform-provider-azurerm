package validate

func LogAnalyticsDataExportWorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	return logAnalyticsGenericName(i, k)
}

func LogAnalyticsDataExportName(i interface{}, k string) (warnings []string, errors []error) {
	return logAnalyticsGenericName(i, k)
}
