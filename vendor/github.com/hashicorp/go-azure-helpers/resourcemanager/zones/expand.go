// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package zones

func Expand(input []string) Schema {
	out := Schema{}

	out = append(out, input...)

	return out
}

func ExpandUntyped(input []interface{}) []string {
	out := make([]string, 0)

	for _, v := range input {
		out = append(out, v.(string))
	}

	return out
}
