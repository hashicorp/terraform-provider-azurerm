// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

// ApplicationSubdomain has the same format requirements as ApplicationName
func ApplicationSubdomain(v interface{}, k string) ([]string, []error) {
	return ApplicationName(v, k)
}
