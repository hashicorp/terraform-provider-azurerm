// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

func CacheNamespacePath(i interface{}, k string) (warnings []string, errs []error) {
	return absolutePath(i, k)
}
