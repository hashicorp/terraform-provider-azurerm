// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validation

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var StringIsNotEmpty schema.SchemaValidateFunc = func(interface{}, string) ([]string, []error) {
	return nil, nil
}

func StringLenBetween(int, int) schema.SchemaValidateFunc {
	return func(interface{}, string) ([]string, []error) {
		return nil, nil
	}
}

func StringInSlice([]string, bool) schema.SchemaValidateFunc {
	return func(interface{}, string) ([]string, []error) {
		return nil, nil
	}
}
