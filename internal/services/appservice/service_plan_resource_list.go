package appservice

import (
	"context"

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

// Metadata implements [sdk.FrameworkListWrappedResource].
func (s *ServicePlanResourceList) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = ServicePlanResource{}.ResourceType()
}

// ResourceFunc implements [sdk.FrameworkListWrappedResource].
func (s *ServicePlanResourceList) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(ServicePlanResource{})
}

// List implements [sdk.FrameworkListWrappedResource].
func (s *ServicePlanResourceList) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
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

	switch {
	case !data.ResourceGroupName.IsNull():
		resourceGroupId := commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString())
		resp, err := client.ListByResourceGroupComplete(ctx, resourceGroupId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, "listing `azurerm_app_service_plan`", err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID), appserviceplans.DefaultListOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, "listing `azurerm_app_service_plan`", err)
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
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing App Service Plan ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, ServicePlanResource{})
			meta.SetID(id)
			if err := resourceServicePlanFlatten(meta, id, &plan); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "flattening App Service Plan data", err)
				return
			}

			sdk.EncodeListResult(ctx, meta.ResourceData, result, push)
		}
	}
}
