package appservice

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/appserviceplans"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServicePlanResourceList struct{}

var _ sdk.FrameworkListWrappedResource = new(ServicePlanResourceList)

func (ServicePlanResourceList) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = ServicePlanResource{}.ResourceType()
}

func (ServicePlanResourceList) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(ServicePlanResource{})
}

func (ServicePlanResourceList) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.AppService.ServicePlanClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []appserviceplans.AppServicePlan

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	r := ServicePlanResource{}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID), appserviceplans.DefaultListOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, plan := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(plan.Name)

			// use case insensitive parsing as the API may return `serverfarms` instead of literal `serverFarms`
			id, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(plan.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing App Service Plan ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)

			if err := r.flatten(meta, id, &plan); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", pointer.From(plan.Name)), err)
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
