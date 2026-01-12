// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

func CacheNFSTargetPath(i interface{}, k string) (warnings []string, errs []error) {
	return relativePath(i, k)
}
