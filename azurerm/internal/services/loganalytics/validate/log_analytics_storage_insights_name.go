package validate

func LogAnalyticsStorageInsightsName(i interface{}, k string) (warnings []string, errors []error) {
	return logAnalyticsGenericName(i, k)
}
