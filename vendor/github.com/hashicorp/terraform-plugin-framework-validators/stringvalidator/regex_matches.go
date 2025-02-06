// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = regexMatchesValidator{}
var _ function.StringParameterValidator = regexMatchesValidator{}

type regexMatchesValidator struct {
	regexp  *regexp.Regexp
	message string
}

func (validator regexMatchesValidator) Description(_ context.Context) string {
	if validator.message != "" {
		return validator.message
	}
	return fmt.Sprintf("value must match regular expression '%s'", validator.regexp)
}

func (validator regexMatchesValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v regexMatchesValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if !v.regexp.MatchString(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

func (v regexMatchesValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	if !v.regexp.MatchString(value) {
		response.Error = validatorfuncerr.InvalidParameterValueMatchFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			value,
		)
	}
}

// RegexMatches returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a string.
//   - Matches the given regular expression https://github.com/google/re2/wiki/Syntax.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
// Optionally an error message can be provided to return something friendlier
// than "value must match regular expression 'regexp'".
func RegexMatches(regexp *regexp.Regexp, message string) regexMatchesValidator {
	return regexMatchesValidator{
		regexp:  regexp,
		message: message,
	}
}
