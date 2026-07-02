package containers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KubernetesFleetAutoUpgradeProfileListResource struct{}

type KubernetesFleetAutoUpgradeProfileListModel struct {
	KubernetesFleetManagerId types.String `tfsdk:"kubernetes_fleet_manager_id"`
}

var _ sdk.FrameworkListWrappedResource = new(KubernetesFleetAutoUpgradeProfileListResource)

func (KubernetesFleetAutoUpgradeProfileListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = KubernetesFleetAutoUpgradeProfileResource{}.ResourceType()
}

func (KubernetesFleetAutoUpgradeProfileListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(KubernetesFleetAutoUpgradeProfileResource{})
}

func (KubernetesFleetAutoUpgradeProfileListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"kubernetes_fleet_manager_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: commonids.ValidateKubernetesFleetID},
				},
			},
		},
	}
}

func (KubernetesFleetAutoUpgradeProfileListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Containers.FleetAutoUpgradeProfilesClient

	var data KubernetesFleetAutoUpgradeProfileListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	fleetId, err := commonids.ParseKubernetesFleetID(data.KubernetesFleetManagerId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Kubernetes Fleet Manager ID for `%s`", KubernetesFleetAutoUpgradeProfileResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ListByFleetComplete(ctx, autoupgradeprofiles.NewFleetID(fleetId.SubscriptionId, fleetId.ResourceGroupName, fleetId.FleetName), autoupgradeprofiles.DefaultListByFleetOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", KubernetesFleetAutoUpgradeProfileResource{}.ResourceType()), err)
		return
	}

	resource := KubernetesFleetAutoUpgradeProfileResource{}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := autoupgradeprofiles.ParseAutoUpgradeProfileID(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Kubernetes Fleet Auto Upgrade Profile ID", err)
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
