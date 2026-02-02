// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package convert

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// convert handles the conversion between Framework types and Go native types used by go-azure-sdk and
// returns reflect.Value for each in the order from -> to and diag.Diagnostics for errors where incorrect types are
// provided. `source` can be a pointer, `target` must be a pointer. It is intended to be used via Expand and Flatten and
// is intentionally not exported.
func convert(source, target any) (reflect.Value, reflect.Value, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	sourceVal := reflect.ValueOf(source)
	targetVal := reflect.ValueOf(target)

	if sourceVal.Kind() == reflect.Ptr {
		sourceVal = sourceVal.Elem()
	}

	if kind := targetVal.Kind(); kind != reflect.Ptr {
		diags.AddError("convert", fmt.Sprintf("target (%T): %s is not a pointer", target, kind))
		return reflect.Value{}, reflect.Value{}, diags
	}

	targetVal = targetVal.Elem()

	return sourceVal, targetVal, diags
}
