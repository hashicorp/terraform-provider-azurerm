// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

func CacheNFSTargetPath(i interface{}, k string) (warnings []string, errs []error) {
	return relativePath(i, k)
}
