// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SubnetListResource struct{}

func (r SubnetListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceSubnet()
}

var _ sdk.FrameworkListWrappedResource = &SubnetListResource{}

type SubnetListModel struct {
	VirtualNetworkID types.String `tfsdk:"virtual_network_id"`
}

func (r SubnetListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = SubnetResourceName
}

func (r SubnetListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"virtual_network_id": listschema.StringAttribute{
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

func (r SubnetListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
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
			// TODO - Do we need to handle limiting the results to ListRequest.Limit?
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(subnet.Name)

			id, err := commonids.ParseSubnetID(pointer.From(subnet.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Virtual Network ID", err)
				return
			}

			vNetResource := resourceSubnet()

			rd := vNetResource.Data(&terraform.InstanceState{})

			rd.SetId(id.ID())

			err = resourceSubnetFlatten(rd, *id, &subnet)
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "encoding Resource data", err)
				return
			}

			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Identity State", err)
				return
			}

			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Identity data", err)
				return
			}

			tfTypeResource, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "converting Resource State data", err)
				return
			}

			if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "setting Resource data", err)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
