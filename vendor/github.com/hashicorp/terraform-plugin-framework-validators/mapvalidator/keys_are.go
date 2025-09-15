// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Map = keysAreValidator{}

// keysAreValidator validates that each map key validates against each of the value validators.
type keysAreValidator struct {
	keyValidators []validator.String
}

// Description describes the validation in plain text formatting.
func (v keysAreValidator) Description(ctx context.Context) string {
	var descriptions []string
	for _, validator := range v.keyValidators {
		descriptions = append(descriptions, validator.Description(ctx))
	}

	return fmt.Sprintf("key must satisfy all validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v keysAreValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateMap performs the validation.
// Note that the Path specified in the MapRequest refers to the value in the Map with key `k`,
// whereas the ConfigValue refers to the key itself (i.e., `k`). This is intentional as the validation being
// performed is for the keys of the Map.
func (v keysAreValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	for k := range req.ConfigValue.Elements() {
		attrPath := req.Path.AtMapKey(k)
		validateReq := validator.StringRequest{
			Path:           attrPath,
			PathExpression: attrPath.Expression(),
			ConfigValue:    types.StringValue(k),
			Config:         req.Config,
		}

		for _, keyValidator := range v.keyValidators {
			validateResp := &validator.StringResponse{}

			keyValidator.ValidateString(ctx, validateReq, validateResp)

			resp.Diagnostics.Append(validateResp.Diagnostics...)
		}
	}
}

// KeysAre returns a map validator that validates all key strings with the
// given string validators.
func KeysAre(keyValidators ...validator.String) validator.Map {
	return keysAreValidator{
		keyValidators: keyValidators,
	}
}
