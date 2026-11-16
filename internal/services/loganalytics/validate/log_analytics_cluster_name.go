// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

func LogAnalyticsClusterName(i any, k string) (warnings []string, errors []error) {
	return logAnalyticsGenericName(i, k)
}
