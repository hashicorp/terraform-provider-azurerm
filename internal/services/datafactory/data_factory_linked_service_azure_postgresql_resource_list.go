// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinkedServiceAzurePostgreSQLListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(LinkedServiceAzurePostgreSQLListResource)

func (LinkedServiceAzurePostgreSQLListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(LinkedServiceAzurePostgreSQLResource{})
}

func (LinkedServiceAzurePostgreSQLListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = LinkedServiceAzurePostgreSQLResource{}.ResourceType()
}

type LinkedServiceAzurePostgreSQLListModel struct {
	DataFactoryID types.String `tfsdk:"data_factory_id"`
}

func (LinkedServiceAzurePostgreSQLListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data_factory_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: linkedservices.ValidateFactoryID,
					},
				},
			},
		},
	}
}

func (LinkedServiceAzurePostgreSQLListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.DataFactory.LinkedServicesClient

	var data LinkedServiceAzurePostgreSQLListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	dataFactoryId, err := linkedservices.ParseFactoryID(data.DataFactoryID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "listing azurerm_data_factory_linked_service_azure_postgresql", err)
		return
	}

	var results []linkedservices.LinkedServiceResource

	r := LinkedServiceAzurePostgreSQLResource{}

	resp, err := client.ListByFactoryComplete(ctx, *dataFactoryId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}

	results = resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, postgreSqlResult := range results {
			if _, ok := postgreSqlResult.Properties.(linkedservices.AzurePostgreSqlLinkedService); !ok {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(postgreSqlResult.Name)

			id, err := linkedservices.ParseLinkedServiceID(pointer.From(postgreSqlResult.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing linked service ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			if err := r.flatten(rmd, id, &postgreSqlResult); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, rmd.ResourceData, &result)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
