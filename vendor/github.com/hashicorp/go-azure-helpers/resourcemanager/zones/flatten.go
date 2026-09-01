// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0

package zones

func Flatten(input *Schema) []string {
	out := make([]string, 0)

	if input != nil {
		out = append(out, *input...)
	}

	return out
}

func FlattenUntyped(input *[]string) []interface{} {
	out := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			out = append(out, v)
		}
	}

	return out
}
