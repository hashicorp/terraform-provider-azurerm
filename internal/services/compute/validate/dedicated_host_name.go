// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

func DedicatedHostName() func(i interface{}, k string) (warnings []string, errors []error) {
	return DedicatedHostGroupName()
}
