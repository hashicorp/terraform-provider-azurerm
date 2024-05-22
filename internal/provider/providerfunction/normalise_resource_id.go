package providerfunction

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

type NormaliseResourceIDFunction struct{}

var _ function.Function = NormaliseResourceIDFunction{}

func NewNormaliseResourceID() function.Function {
	return &NormaliseResourceIDFunction{}
}

func (a NormaliseResourceIDFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "normalise_resource_id"
}

func (a NormaliseResourceIDFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary:             "normalise_resource_id",
		Description:         "Parses and normalises the casing on an Azure Resource Manager ID into the correct casing for Terraform",
		MarkdownDescription: "Parses and normalises the casing on an Azure Resource Manager ID into the correct casing for Terraform",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "id",
				Description:         "Resource ID",
				MarkdownDescription: "Resource ID",
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

	response.Error = function.ConcatFuncErrors(response.Result.Set(ctx, recaser.ReCase(id)))
}
