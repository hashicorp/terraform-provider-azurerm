// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subnet

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ListResource struct{}

func (r ListResource) ResourceFunc() *pluginsdk.Resource {
	return RegisterResource()
}

var _ sdk.FrameworkListWrappedResource = &ListResource{}

type SubnetListModel struct {
	VirtualNetworkID types.String `tfsdk:"virtual_network_id"`
}

func (r ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = SubnetResourceName
}

func (r ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"virtual_network_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateVirtualNetworkID,
					},
				},
			},
		},
	}
}

func (r ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Network.Subnets

	var data SubnetListModel

	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	vnetID, err := commonids.ParseVirtualNetworkID(data.VirtualNetworkID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Virtual Network ID", err)
	}

	resp, err := client.ListComplete(ctx, *vnetID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing %s", VirtualNetworkResourceName), err)
		return
	}

	listResults := resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, subnet := range listResults {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(subnet.Name)

			id, err := commonids.ParseSubnetID(pointer.From(subnet.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Virtual Network ID", err)
				return
			}

			rd := RegisterResource().Data(&terraform.InstanceState{})

			rd.SetId(id.ID())

			if err = resourceSubnetFlatten(rd, *id, &subnet); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "encoding resource data", err)
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
