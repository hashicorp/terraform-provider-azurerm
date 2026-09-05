// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	workbooktemplates "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApplicationInsightsWorkbookTemplateListResource struct{}

type ApplicationInsightsWorkbookTemplateListModel struct {
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
}

var _ sdk.FrameworkListWrappedResourceWithConfig = new(ApplicationInsightsWorkbookTemplateListResource)

func (ApplicationInsightsWorkbookTemplateListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = ApplicationInsightsWorkbookTemplateResource{}.ResourceType()
}

func (ApplicationInsightsWorkbookTemplateListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(ApplicationInsightsWorkbookTemplateResource{})
}

func (ApplicationInsightsWorkbookTemplateListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"resource_group_name": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: resourcegroups.ValidateName},
				},
			},
		},
	}
}

func (ApplicationInsightsWorkbookTemplateListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.AppInsights.WorkbookTemplateClient

	var data ApplicationInsightsWorkbookTemplateListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	resource := ApplicationInsightsWorkbookTemplateResource{}

	resp, err := client.WorkbookTemplatesListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(metadata.SubscriptionId, data.ResourceGroupName.ValueString()))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", resource.ResourceType()), err)
		return
	}

	results := resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := workbooktemplates.ParseWorkbookTemplateIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing %s ID", resource.ResourceType()), err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, resource)
			meta.SetID(id)

			if err := resource.flatten(meta, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", resource.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, meta.ResourceData, &result)
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
