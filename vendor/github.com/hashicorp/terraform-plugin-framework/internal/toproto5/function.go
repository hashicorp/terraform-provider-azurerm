// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// Function returns the *tfprotov5.Function for a function.Definition.
func Function(ctx context.Context, fw function.Definition) *tfprotov5.Function {
	proto := &tfprotov5.Function{
		DeprecationMessage: fw.DeprecationMessage,
		Parameters:         make([]*tfprotov5.FunctionParameter, 0, len(fw.Parameters)),
		Return:             FunctionReturn(ctx, fw.Return),
		Summary:            fw.Summary,
	}

	if fw.MarkdownDescription != "" {
		proto.Description = fw.MarkdownDescription
		proto.DescriptionKind = tfprotov5.StringKindMarkdown
	} else if fw.Description != "" {
		proto.Description = fw.Description
		proto.DescriptionKind = tfprotov5.StringKindPlain
	}

	for _, fwParameter := range fw.Parameters {
		protoParam := FunctionParameter(ctx, fwParameter)
		proto.Parameters = append(proto.Parameters, protoParam)
	}

	if fw.VariadicParameter != nil {
		protoParam := FunctionParameter(ctx, fw.VariadicParameter)
		proto.VariadicParameter = protoParam
	}

	return proto
}

// FunctionParameter returns the *tfprotov5.FunctionParameter for a
// function.Parameter.
func FunctionParameter(ctx context.Context, fw function.Parameter) *tfprotov5.FunctionParameter {
	if fw == nil {
		return nil
	}

	proto := &tfprotov5.FunctionParameter{
		AllowNullValue:     fw.GetAllowNullValue(),
		AllowUnknownValues: fw.GetAllowUnknownValues(),
		Name:               fw.GetName(),
		Type:               fw.GetType().TerraformType(ctx),
	}

	if fw.GetMarkdownDescription() != "" {
		proto.Description = fw.GetMarkdownDescription()
		proto.DescriptionKind = tfprotov5.StringKindMarkdown
	} else if fw.GetDescription() != "" {
		proto.Description = fw.GetDescription()
		proto.DescriptionKind = tfprotov5.StringKindPlain
	}

	return proto
}

// FunctionMetadata returns the tfprotov5.FunctionMetadata for a
// fwserver.FunctionMetadata.
func FunctionMetadata(ctx context.Context, fw fwserver.FunctionMetadata) tfprotov5.FunctionMetadata {
	proto := tfprotov5.FunctionMetadata{
		Name: fw.Name,
	}

	return proto
}

// FunctionReturn returns the *tfprotov5.FunctionReturn for a
// function.Return.
func FunctionReturn(ctx context.Context, fw function.Return) *tfprotov5.FunctionReturn {
	if fw == nil {
		return nil
	}

	proto := &tfprotov5.FunctionReturn{
		Type: fw.GetType().TerraformType(ctx),
	}

	return proto
}

// FunctionResultData returns the *tfprotov5.DynamicValue for a given
// function.ResultData.
func FunctionResultData(ctx context.Context, data function.ResultData) (*tfprotov5.DynamicValue, *function.FuncError) {
	attrValue := data.Value()

	if attrValue == nil {
		return nil, nil
	}

	tfType := attrValue.Type(ctx).TerraformType(ctx)
	tfValue, err := attrValue.ToTerraformValue(ctx)

	if err != nil {
		msg := "Unable to Convert Function Result Data: An unexpected error was encountered when converting the function result data to the protocol type. " +
			"Please report this to the provider developer:\n\n" +
			"Unable to convert framework type to tftypes: " + err.Error()

		return nil, function.NewFuncError(msg)
	}

	dynamicValue, err := tfprotov5.NewDynamicValue(tfType, tfValue)

	if err != nil {
		msg := "Unable to Convert Function Result Data: An unexpected error was encountered when converting the function result data to the protocol type. " +
			"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n" +
			"Unable to create DynamicValue: " + err.Error()

		return nil, function.NewFuncError(msg)
	}

	return &dynamicValue, nil
}
