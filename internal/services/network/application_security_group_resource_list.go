// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/applicationsecuritygroups"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApplicationSecurityGroupListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(ApplicationSecurityGroupListResource)

func (r ApplicationSecurityGroupListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceApplicationSecurityGroup()
}

func (r ApplicationSecurityGroupListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_application_security_group"
}

func (r ApplicationSecurityGroupListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.ApplicationSecurityGroups
	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]applicationsecuritygroups.ApplicationSecurityGroup, 0)
	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_application_security_group`), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListAllComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", `azurerm_application_security_group`), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, applicationsecuritygroup := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(applicationsecuritygroup.Name)

			rd := resourceApplicationSecurityGroup().Data(&terraform.InstanceState{})

			id, err := applicationsecuritygroups.ParseApplicationSecurityGroupID(pointer.From(applicationsecuritygroup.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Network ApplicationSecurityGroup ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourceApplicationSecurityGroupFlatten(rd, id, &applicationsecuritygroup); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_application_security_group"), err)
				return
			}

			sdk.EncodeListResult(ctx, rd, &result)
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
