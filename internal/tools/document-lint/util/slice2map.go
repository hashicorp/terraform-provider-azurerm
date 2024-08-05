// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

func Slice2Map(in []string) (out map[string]struct{}) {
	out = make(map[string]struct{}, len(in))
	for _, k := range in {
		out[k] = struct{}{}
	}
	return
}
