// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// All returns a validator which ensures that any configured attribute value
// attribute value validates against all the given validators.
//
// Use of All is only necessary when used in conjunction with Any or AnyWithAllWarnings
// as the Validators field automatically applies a logical AND.
func All(validators ...validator.Map) validator.Map {
	return allValidator{
		validators: validators,
	}
}

var _ validator.Map = allValidator{}

// allValidator implements the validator.
type allValidator struct {
	validators []validator.Map
}

// Description describes the validation in plain text formatting.
func (v allValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, subValidator := range v.validators {
		descriptions = append(descriptions, subValidator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy all of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v allValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateMap performs the validation.
func (v allValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	for _, subValidator := range v.validators {
		validateResp := &validator.MapResponse{}

		subValidator.ValidateMap(ctx, req, validateResp)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}
