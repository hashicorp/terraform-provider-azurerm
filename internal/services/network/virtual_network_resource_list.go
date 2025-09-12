// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type VirtualNetworkListResource struct {
	sdk.ListResourceMetadata
}

var _ sdk.ListResourceWithRawV5Schemas = &VirtualNetworkListResource{}

type VirtualNetworkListModel struct {
	ResourceGroupName string `tfsdk:"resource_group_name"`
}

func NewVirtualNetworkListResource() list.ListResource {
	return &VirtualNetworkListResource{}
}

func (r *VirtualNetworkListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = VirtualNetworkResourceName
}

func (r *VirtualNetworkListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, response *list.RawV5SchemaResponse) {
	res := resourceVirtualNetwork()
	response.ProtoV5Schema = res.ProtoSchema(ctx)()
	response.ProtoV5IdentitySchema = res.ProtoIdentitySchema(ctx)()
}

func (r *VirtualNetworkListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"resource_group_name": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
				},
			},
		},
	}
}

func (r *VirtualNetworkListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream) {
	client := r.Client.Network.VirtualNetworks
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute) // TODO - Is this long enough for a list call?
	defer cancel()

	var data VirtualNetworkListModel

	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	resp, err := client.ListComplete(ctx, commonids.NewResourceGroupID(r.SubscriptionId, data.ResourceGroupName))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing %s", VirtualNetworkResourceName), err)
		return
	}

	listResults := resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, vnet := range listResults {
			// TODO - Do we need to handle limiting the results to ListRequest.Limit?
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(vnet.Name)

			id, err := commonids.ParseVirtualNetworkID(*vnet.Id)
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "parsing Virtual Network ID", err)
				return
			}

			vNetResource := resourceVirtualNetwork()

			rd := vNetResource.Data(&terraform.InstanceState{})

			rd.SetId(id.ID())

			err = resourceVirtualNetworkFlatten(rd, *id, &vnet)
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "encoding resource data", err)
				return
			}

			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "converting Identity State", err)
				return
			}

			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "setting identity data", err)
				return
			}

			tfTypeResource, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "converting Resource State data", err)
				return
			}

			if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "setting resource data", err)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}

func (r *VirtualNetworkListResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Defaults(request, response)
}
