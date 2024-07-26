// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frameworkhelpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// WrappedStringValidator provides a wrapper for legacy SDKv2 type validations to ease migration to Framework Native
// The provided function is tested against the value in the configuration and populates the diagnostics accordingly.
type WrappedStringValidator struct {
	Func         func(v interface{}, k string) (warnings []string, errors []error)
	Desc         string
	MarkdownDesc string
}

func (w WrappedStringValidator) Description(_ context.Context) string {
	return w.Desc
}

func (w WrappedStringValidator) MarkdownDescription(_ context.Context) string {
	return w.MarkdownDesc
}

func (w WrappedStringValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	path := request.Path.String()
	warnings, err := w.Func(value, path)

	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("invalid value for %s", path), fmt.Sprintf("%+v", err))
		return
	}

	if len(warnings) > 0 { // This may be redundant - legacy validators never really used warnings.
		for _, v := range warnings {
			response.Diagnostics.Append(diag.NewWarningDiagnostic(fmt.Sprintf("validating %s", path), v))
		}
	}
}

var _ validator.String = &WrappedStringValidator{}

type WrappedStringDefault struct {
	Desc     *string
	Markdown *string
	Value    string
}

var _ defaults.String = WrappedStringDefault{}

// NewWrappedStringDefault is a helper function to return a new defaults.String implementation for any type that
// implements the Go string type.
func NewWrappedStringDefault[T ~string](value T) WrappedStringDefault {
	return WrappedStringDefault{
		Value: string(value),
	}
}

func (w WrappedStringDefault) Description(_ context.Context) string {
	return pointer.From(w.Desc)
}

func (w WrappedStringDefault) MarkdownDescription(_ context.Context) string {
	return pointer.From(w.Markdown)
}

func (w WrappedStringDefault) DefaultString(_ context.Context, _ defaults.StringRequest, response *defaults.StringResponse) {
	d := basetypes.NewStringValue(w.Value)
	response.PlanValue = d
}
