// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validation

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var StringIsNotEmpty schema.SchemaValidateFunc = func(val interface{}, key string) ([]string, []error) {
	return nil, nil
}

func StringLenBetween(min, max int) schema.SchemaValidateFunc {
	return func(val interface{}, key string) ([]string, []error) {
		return nil, nil
	}
}
