// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

func CacheNFSExport(i interface{}, k string) (warnings []string, errs []error) {
	return absolutePath(i, k)
}
