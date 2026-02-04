// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package workloads

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WorkloadsSAPSingleNodeVirtualInstanceListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(WorkloadsSAPSingleNodeVirtualInstanceListResource)

func (r WorkloadsSAPSingleNodeVirtualInstanceListResource) ResourceFunc() *pluginsdk.Resource {
	// Wrap the SDK v2 typed resource and convert it to pluginsdk.Resource
	wrapper := sdk.NewResourceWrapper(WorkloadsSAPSingleNodeVirtualInstanceResource{})
	resource, err := wrapper.Resource()
	if err != nil {
		panic(fmt.Sprintf("failed to wrap resource: %+v", err))
	}
	return resource
}

func (r WorkloadsSAPSingleNodeVirtualInstanceListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_workloads_sap_single_node_virtual_instance"
}

func (r WorkloadsSAPSingleNodeVirtualInstanceListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Workloads.SAPVirtualInstances

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]sapvirtualinstances.SAPVirtualInstance, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_workloads_sap_single_node_virtual_instance"), err)
			return
		}
		results = resp.Items
	default:
		resp, err := client.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_workloads_sap_single_node_virtual_instance"), err)
			return
		}
		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, instance := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(instance.Name)

			wrapper := sdk.NewResourceWrapper(WorkloadsSAPSingleNodeVirtualInstanceResource{})
			resource, err := wrapper.Resource()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "creating resource wrapper", err)
				return
			}
			rd := resource.Data(&terraform.InstanceState{})

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(pointer.From(instance.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing SAP Virtual Instance ID", err)
				return
			}
			rd.SetId(id.ID())

			if err := flattenSingleNodeForListResource(rd, id, &instance); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("flattening `%s` resource data", "azurerm_workloads_sap_single_node_virtual_instance"), err)
				return
			}

			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Identity State", err)
				return
			}
			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Identity Data", err)
				return
			}

			tfTypeResourceState, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Resource State", err)
				return
			}
			if err := result.Resource.Set(ctx, *tfTypeResourceState); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Resource Data", err)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
