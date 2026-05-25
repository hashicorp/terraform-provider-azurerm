// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

func LogAnalyticsDataExportName(i interface{}, k string) (warnings []string, errors []error) {
	return logAnalyticsGenericName(i, k)
}
