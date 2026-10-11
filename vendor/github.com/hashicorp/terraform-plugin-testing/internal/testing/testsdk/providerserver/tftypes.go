// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func DynamicValueToValue(schema *tfprotov6.Schema, dynamicValue *tfprotov6.DynamicValue) (tftypes.Value, *tfprotov6.Diagnostic) {
	if schema == nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Convert DynamicValue",
			Detail:   "Converting the DynamicValue to Value returned an unexpected error: missing schema",
		}

		return tftypes.NewValue(tftypes.Object{}, nil), diag
	}

	if dynamicValue == nil {
		return tftypes.NewValue(schema.ValueType(), nil), nil
	}

	value, err := dynamicValue.Unmarshal(schema.ValueType())

	if err != nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Convert DynamicValue",
			Detail:   "Converting the DynamicValue to Value returned an unexpected error: " + err.Error(),
		}

		return value, diag
	}

	return value, nil
}

func ValuetoDynamicValue(schema *tfprotov6.Schema, value tftypes.Value) (*tfprotov6.DynamicValue, *tfprotov6.Diagnostic) {
	if schema == nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Convert Value",
			Detail:   "Converting the Value to DynamicValue returned an unexpected error: missing schema",
		}

		return nil, diag
	}

	dynamicValue, err := tfprotov6.NewDynamicValue(schema.ValueType(), value)

	if err != nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Convert Value",
			Detail:   "Converting the Value to DynamicValue returned an unexpected error: " + err.Error(),
		}

		return &dynamicValue, diag
	}

	return &dynamicValue, nil
}

func IdentityDynamicValueToValue(schema *tfprotov6.ResourceIdentitySchema, dynamicValue *tfprotov6.DynamicValue) (tftypes.Value, *tfprotov6.Diagnostic) {
	if schema == nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Convert DynamicValue",
			Detail:   "Converting the DynamicValue to Value returned an unexpected error: missing identity schema",
		}

		return tftypes.NewValue(tftypes.Object{}, nil), diag
	}

	if dynamicValue == nil {
		return tftypes.NewValue(schema.ValueType(), nil), nil
	}

	value, err := dynamicValue.Unmarshal(schema.ValueType())

	if err != nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Convert DynamicValue",
			Detail:   "Converting the DynamicValue to Value returned an unexpected error: " + err.Error(),
		}

		return value, diag
	}

	return value, nil
}

func IdentityValuetoDynamicValue(schema *tfprotov6.ResourceIdentitySchema, value tftypes.Value) (*tfprotov6.DynamicValue, *tfprotov6.Diagnostic) {
	if schema == nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Convert Value",
			Detail:   "Converting the Value to DynamicValue returned an unexpected error: missing identity schema",
		}

		return nil, diag
	}

	dynamicValue, err := tfprotov6.NewDynamicValue(schema.ValueType(), value)

	if err != nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Convert Value",
			Detail:   "Converting the Value to DynamicValue returned an unexpected error: " + err.Error(),
		}

		return &dynamicValue, diag
	}

	return &dynamicValue, nil
}
