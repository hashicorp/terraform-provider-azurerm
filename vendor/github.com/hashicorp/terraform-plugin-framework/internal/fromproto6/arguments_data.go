// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ArgumentsData returns the ArgumentsData for a given []*tfprotov6.DynamicValue
// and function.Definition.
func ArgumentsData(ctx context.Context, arguments []*tfprotov6.DynamicValue, definition function.Definition) (function.ArgumentsData, diag.Diagnostics) {
	if definition.VariadicParameter == nil && len(arguments) != len(definition.Parameters) {
		return function.NewArgumentsData(nil), diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Unexpected Function Arguments Data",
				"The provider received an unexpected number of function arguments from Terraform for the given function definition. "+
					"This is always an issue in terraform-plugin-framework or Terraform itself and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Expected function arguments: %d\n", len(definition.Parameters))+
					fmt.Sprintf("Given function arguments: %d", len(arguments)),
			),
		}
	}

	// Expect at least all parameters to have corresponding arguments. Variadic
	// parameter might have 0 to n arguments, which is why it is not checked in
	// this case.
	if len(arguments) < len(definition.Parameters) {
		return function.NewArgumentsData(nil), diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Unexpected Function Arguments Data",
				"The provider received an unexpected number of function arguments from Terraform for the given function definition. "+
					"This is always an issue in terraform-plugin-framework or Terraform itself and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Expected minimum function arguments: %d\n", len(definition.Parameters))+
					fmt.Sprintf("Given function arguments: %d", len(arguments)),
			),
		}
	}

	if definition.VariadicParameter == nil && len(arguments) == 0 {
		return function.NewArgumentsData(nil), nil
	}

	// Variadic values are collected as a separate tuple to ease developer usage.
	argumentValues := make([]attr.Value, 0, len(definition.Parameters))
	variadicValues := make([]attr.Value, 0, len(arguments)-len(definition.Parameters))
	var diags diag.Diagnostics

	for position, argument := range arguments {
		parameter, parameterDiags := definition.Parameter(ctx, position)

		diags.Append(parameterDiags...)

		if diags.HasError() {
			return function.NewArgumentsData(nil), diags
		}

		parameterType := parameter.GetType()

		if parameterType == nil {
			diags.AddError(
				"Unable to Convert Function Argument",
				"An unexpected error was encountered when converting the function argument from the protocol type. "+
					"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
					"Please report this to the provider developer:\n\n"+
					fmt.Sprintf("Parameter type missing at position %d", position),
			)

			return function.NewArgumentsData(nil), diags
		}

		tfValue, err := argument.Unmarshal(parameterType.TerraformType(ctx))

		if err != nil {
			diags.AddError(
				"Unable to Convert Function Argument",
				"An unexpected error was encountered when converting the function argument from the protocol type. "+
					"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
					"Please report this to the provider developer:\n\n"+
					fmt.Sprintf("Unable to unmarshal DynamicValue at position %d: %s", position, err),
			)

			return function.NewArgumentsData(nil), diags
		}

		attrValue, err := parameterType.ValueFromTerraform(ctx, tfValue)

		if err != nil {
			diags.AddError(
				"Unable to Convert Function Argument",
				"An unexpected error was encountered when converting the function argument from the protocol type. "+
					"Please report this to the provider developer:\n\n"+
					fmt.Sprintf("Unable to convert tftypes to framework type at position %d: %s", position, err),
			)

			return function.NewArgumentsData(nil), diags
		}

		// This is intentionally below the attr.Value conversion so it can be
		// updated for any new type system validation interfaces. Note that the
		// original xattr.TypeWithValidation interface must set a path.Path,
		// which will always be incorrect in the context of functions.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/589
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/893
		if attrTypeWithValidate, ok := parameterType.(xattr.TypeWithValidate); ok {
			logging.FrameworkTrace(ctx, "Parameter type implements TypeWithValidate")
			logging.FrameworkTrace(ctx, "Calling provider defined Type Validate")
			diags.Append(attrTypeWithValidate.Validate(ctx, tfValue, path.Empty())...)
			logging.FrameworkTrace(ctx, "Called provider defined Type Validate")

			if diags.HasError() {
				continue
			}
		}

		if definition.VariadicParameter != nil && position >= len(definition.Parameters) {
			variadicValues = append(variadicValues, attrValue)

			continue
		}

		argumentValues = append(argumentValues, attrValue)
	}

	if definition.VariadicParameter != nil {
		// MAINTAINER NOTE: Variadic parameters are represented as individual arguments in the CallFunction RPC and Terraform core applies the variadic parameter
		// type constraint to each argument individually. For developer convenience, the framework logic below, groups the variadic arguments into a
		// framework Tuple where each element type of the tuple matches the variadic parameter type.
		//
		// Previously, this logic utilized a framework List with an element type that matched the variadic parameter type. Using a List presented an issue with dynamic
		// variadic parameters, as each argument was allowed to be any type "individually", rather than having a single type constraint applied to all dynamic elements,
		// like a cty.List in Terraform. This eventually results in an error attempting to create a tftypes.List with multiple element types (when unwrapping from a framework
		// dynamic to a tftypes concrete value).
		//
		// While a framework List type can handle multiple dynamic values of different types (due to it's wrapping of dynamic values), `terraform-plugin-go` and `tftypes.List` cannot.
		// Currently, the logic for retrieving argument data is dependent on the tftypes package to utilize the framework reflection logic, requiring us to apply a type constraint
		// that is valid in Terraform core and `terraform-plugin-go`, which we are doing here with a Tuple.
		variadicType := definition.VariadicParameter.GetType()
		tupleTypes := make([]attr.Type, len(variadicValues))
		tupleValues := make([]attr.Value, len(variadicValues))
		for i, val := range variadicValues {
			tupleTypes[i] = variadicType
			tupleValues[i] = val
		}
		variadicValue, variadicValueDiags := basetypes.NewTupleValue(tupleTypes, tupleValues)

		diags.Append(variadicValueDiags...)

		if diags.HasError() {
			return function.NewArgumentsData(argumentValues), diags
		}

		argumentValues = append(argumentValues, variadicValue)
	}

	return function.NewArgumentsData(argumentValues), diags
}
