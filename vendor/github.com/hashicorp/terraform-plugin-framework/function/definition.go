// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwfunction"
)

// Definition is a function definition. Always set at least the Result field.
type Definition struct {
	// Parameters is the ordered list of function parameters and their
	// associated data types.
	Parameters []Parameter

	// VariadicParameter is an optional final parameter which can accept zero or
	// more arguments when the function is called. The argument data is sent as
	// a tuple, where all elements are of the same associated data type.
	VariadicParameter Parameter

	// Return is the function call response data type.
	Return Return

	// Summary is a short description of the function, preferably a single
	// sentence. Use the Description field for longer documentation about the
	// function and its implementation.
	Summary string

	// Description is the longer documentation for usage, such as editor
	// integrations, to give practitioners more information about the purpose of
	// the function and how its logic is implemented. It should be plaintext
	// formatted.
	Description string

	// MarkdownDescription is the longer documentation for usage, such as a
	// registry, to give practitioners more information about the purpose of the
	// function and how its logic is implemented.
	MarkdownDescription string

	// DeprecationMessage defines warning diagnostic details to display when
	// practitioner configurations use this function. The warning diagnostic
	// summary is automatically set to "Function Deprecated" along with
	// configuration source file and line information.
	DeprecationMessage string
}

// ValidateImplementation contains logic for validating the provider-defined
// implementation of the definition to prevent unexpected errors or panics. This
// logic runs during the GetProviderSchema RPC, or via provider-defined unit
// testing, and should never include false positives.
func (d Definition) ValidateImplementation(ctx context.Context, req DefinitionValidateRequest, resp *DefinitionValidateResponse) {
	var diags diag.Diagnostics

	if d.Return == nil {
		diags.AddError(
			"Invalid Function Definition",
			"When validating the function definition, an implementation issue was found. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				fmt.Sprintf("Function %q - Definition Return field is undefined", req.FuncName),
		)
	} else if d.Return.GetType() == nil {
		diags.AddError(
			"Invalid Function Definition",
			"When validating the function definition, an implementation issue was found. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				fmt.Sprintf("Function %q - Definition return data type is undefined", req.FuncName),
		)
	} else if returnWithValidateImplementation, ok := d.Return.(fwfunction.ReturnWithValidateImplementation); ok {
		req := fwfunction.ValidateReturnImplementationRequest{}
		resp := &fwfunction.ValidateReturnImplementationResponse{}

		returnWithValidateImplementation.ValidateImplementation(ctx, req, resp)

		diags.Append(resp.Diagnostics...)
	}

	paramNames := make(map[string]int, len(d.Parameters))
	for pos, param := range d.Parameters {
		parameterPosition := int64(pos)
		name := param.GetName()
		// If name is not set, add an error diagnostic, parameter names are mandatory.
		if name == "" {
			diags.AddError(
				"Invalid Function Definition",
				"When validating the function definition, an implementation issue was found. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Function %q - Parameter at position %d does not have a name", req.FuncName, pos),
			)
		}

		if paramWithValidateImplementation, ok := param.(fwfunction.ParameterWithValidateImplementation); ok {
			req := fwfunction.ValidateParameterImplementationRequest{
				Name:              name,
				ParameterPosition: &parameterPosition,
			}
			resp := &fwfunction.ValidateParameterImplementationResponse{}

			paramWithValidateImplementation.ValidateImplementation(ctx, req, resp)

			diags.Append(resp.Diagnostics...)
		}

		conflictPos, exists := paramNames[name]
		if exists && name != "" {
			diags.AddError(
				"Invalid Function Definition",
				"When validating the function definition, an implementation issue was found. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					"Parameter names must be unique. "+
					fmt.Sprintf("Function %q - Parameters at position %d and %d have the same name %q", req.FuncName, conflictPos, pos, name),
			)
			continue
		}

		paramNames[name] = pos
	}

	if d.VariadicParameter != nil {
		name := d.VariadicParameter.GetName()
		// If name is not set, add an error diagnostic, parameter names are mandatory.
		if name == "" {
			diags.AddError(
				"Invalid Function Definition",
				"When validating the function definition, an implementation issue was found. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Function %q - The variadic parameter does not have a name", req.FuncName),
			)
		}

		if paramWithValidateImplementation, ok := d.VariadicParameter.(fwfunction.ParameterWithValidateImplementation); ok {
			req := fwfunction.ValidateParameterImplementationRequest{
				Name: name,
			}
			resp := &fwfunction.ValidateParameterImplementationResponse{}

			paramWithValidateImplementation.ValidateImplementation(ctx, req, resp)

			diags.Append(resp.Diagnostics...)
		}

		conflictPos, exists := paramNames[name]
		if exists && name != "" {
			diags.AddError(
				"Invalid Function Definition",
				"When validating the function definition, an implementation issue was found. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					"Parameter names must be unique. "+
					fmt.Sprintf("Function %q - Parameter at position %d and the variadic parameter have the same name %q", req.FuncName, conflictPos, name),
			)
		}
	}

	resp.Diagnostics.Append(diags...)
}

// DefinitionRequest represents a request for the Function to return its
// definition, such as its ordered parameters and result. An instance of this
// request struct is supplied as an argument to the Function type Definition
// method.
type DefinitionRequest struct{}

// DefinitionResponse represents a response to a DefinitionRequest. An instance
// of this response struct is supplied as an argument to the Function type
// Definition method. Always set at least the Definition field.
type DefinitionResponse struct {
	// Definition is the function definition.
	Definition Definition

	// Diagnostics report errors or warnings related to defining the function.
	// An empty slice indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics
}

// DefinitionValidateRequest represents a request for the Function to validate its
// definition. An instance of this request struct is supplied as an argument to
// the Definition type ValidateImplementation method.
type DefinitionValidateRequest struct {
	// FuncName is the name of the function definition being validated.
	FuncName string
}

// DefinitionValidateResponse represents a response to a DefinitionValidateRequest.
// An instance of this response struct is supplied as an argument to the Definition
// type ValidateImplementation method.
type DefinitionValidateResponse struct {
	// Diagnostics report errors or warnings related to validation of a function
	// definition. An empty slice indicates success, with no warnings or errors
	// generated.
	Diagnostics diag.Diagnostics
}
