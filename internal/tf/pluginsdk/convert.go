// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/gocty"
)

// GoValueFromTerraformValue returns a pointer to the Native Go value for the provided input cty.Value
// If the input value is null, a nil pointer for the type is returned.
// This is a generics function, usage requires supplying the expected type.
// e.g. out, err := GoValueFromTerraformValue[string](someVal)
// NOTE: This helper is experimental and should only be used by Hashicorp maintainers until further notice
func GoValueFromTerraformValue[T any](input cty.Value) (*T, error) {
	result := new(T)
	if input.IsNull() {
		return result, nil
	}
	if err := gocty.FromCtyValue(input, result); err != nil {
		return nil, err
	}

	return result, nil
}
