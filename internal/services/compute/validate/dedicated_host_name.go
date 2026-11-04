// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

func DedicatedHostName() func(i any, k string) (warnings []string, errors []error) {
	return DedicatedHostGroupName()
}
