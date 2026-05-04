// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// Return is the interface for defining function return data.
type Return interface {
	// GetType should return the data type for the return, which determines
	// what data type Terraform requires for configurations receiving the
	// response of a function call and the return data type required from the
	// Function type Run method.
	GetType() attr.Type

	// NewResultData should return a new ResultData with an unknown value (or
	// best approximation of an invalid value) of the corresponding data type.
	// The Function type Run method is expected to overwrite the value before
	// returning.
	NewResultData(context.Context) (ResultData, *FuncError)
}
