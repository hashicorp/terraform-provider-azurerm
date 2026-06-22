// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managedidentity

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type UserAssignedIdentityListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(UserAssignedIdentityListResource)

func (UserAssignedIdentityListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = UserAssignedIdentityResource{}.ResourceType()
}

func (UserAssignedIdentityListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(UserAssignedIdentityResource{})
}

func (UserAssignedIdentityListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ManagedIdentity.V20241130.Identities

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	r := UserAssignedIdentityResource{}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.UserAssignedIdentitiesListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
			return
		}

		stream.Results = func(push func(list.ListResult) bool) {
			for _, item := range resp.Items {
				result := request.NewListResult(ctx)
				result.DisplayName = pointer.From(item.Name)

				id, err := commonids.ParseUserAssignedIdentityIDInsensitively(pointer.From(item.Id))
				if err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing User Assigned Identity ID", err)
					return
				}

				meta := sdk.NewResourceMetaData(metadata.Client, r)
				meta.SetID(id)

				if err := r.flatten(meta, id, &item); err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
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
	default:
		resp, err := client.UserAssignedIdentitiesListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
			return
		}

		stream.Results = func(push func(list.ListResult) bool) {
			for _, item := range resp.Items {
				result := request.NewListResult(ctx)
				result.DisplayName = pointer.From(item.Name)

				id, err := commonids.ParseUserAssignedIdentityIDInsensitively(pointer.From(item.Id))
				if err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing User Assigned Identity ID", err)
					return
				}

				meta := sdk.NewResourceMetaData(metadata.Client, r)
				meta.SetID(id)

				if err := r.flatten(meta, id, &item); err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
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
}
