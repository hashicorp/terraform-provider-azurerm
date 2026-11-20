// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

type NormaliseResourceIDFunction struct{}

var _ function.Function = NormaliseResourceIDFunction{}

func NewNormaliseResourceIDFunction() function.Function {
	return &NormaliseResourceIDFunction{}
}

func (a NormaliseResourceIDFunction) Metadata(_ context.Context, _ function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "normalise_resource_id"
}

func (a NormaliseResourceIDFunction) Definition(_ context.Context, _ function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary: "normalise_resource_id",
		Description: "Takes an Azure Resource ID and attempts to normalise the case-sensitive system segments as required by the AzureRM provider.\n\n" +
			"~> **Note:** User specified segments are not affected or corrected. (e.g. resource names). Please ensure that these match your configuration correctly to avoid errors. If a resource is not supported by the provider, this function may not provide a correct result.",
		MarkdownDescription: "Takes an Azure Resource ID and attempts to normalise the case-sensitive system segments as required by the AzureRM provider.\n\n" +
			"~> **Note:** User specified segments are not affected or corrected. (e.g. resource names). Please ensure that these match your configuration correctly to avoid errors. If a resource is not supported by the provider, this function may not provide a correct result.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "id",
				Description:         "Azure Resource Manager ID.",
				MarkdownDescription: "Azure Resource Manager ID.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (a NormaliseResourceIDFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var id string

	response.Error = function.ConcatFuncErrors(request.Arguments.Get(ctx, &id))

	if response.Error != nil {
		return
	}

	result, err := recaser.ReCaseKnownId(id)
	if err != nil {
		response.Error = function.NewFuncError(err.Error())
		return
	}

	response.Error = function.ConcatFuncErrors(response.Result.Set(ctx, result))
}
