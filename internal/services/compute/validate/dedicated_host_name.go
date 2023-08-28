// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

func DedicatedHostName() func(i interface{}, k string) (warnings []string, errors []error) {
	return DedicatedHostGroupName()
}
