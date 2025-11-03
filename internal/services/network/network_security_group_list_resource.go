package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecuritygroups"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetworkSecurityGroupListResource struct {
	sdk.ListResourceMetadata
}

type NetworkSecurityGroupListModel struct {
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	SubscriptionId    types.String `tfsdk:"subscription_id"`
}

var _ sdk.ListResourceWithRawV5Schemas = new(NetworkSecurityGroupListResource)

func NewNetworkSecurityGroupListResource() list.ListResource {
	return new(NetworkSecurityGroupListResource)
}

func (r *NetworkSecurityGroupListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = networkSecurityGroupResourceName
}

func (r *NetworkSecurityGroupListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_group_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
				},
			},
			"subscription_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.IsUUID,
					},
				},
			},
		},
	}
}

func (r *NetworkSecurityGroupListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream) {
	client := r.Client.Network.NetworkSecurityGroups

	ctx, cancel := context.WithTimeout(ctx, time.Minute*60)
	defer cancel()

	var data NetworkSecurityGroupListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]networksecuritygroups.NetworkSecurityGroup, 0)

	subscriptionID := r.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", networkSecurityGroupResourceName), err)
			return
		}

		results = resp.Items
	default:
		resp, err := client.ListAllComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", networkSecurityGroupResourceName), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, nsg := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(nsg.Name)

			id, err := networksecuritygroups.ParseNetworkSecurityGroupID(pointer.From(nsg.Id))
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "parsing Network Security Group ID", err)
				return
			}

			rd := resourceNetworkSecurityGroup().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceNetworkSecurityGroupFlatten(rd, id, &nsg); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("encoding `%s` resource data", networkSecurityGroupResourceName), err)
				return
			}

			tfTypeIdentity, err := rd.TfTypeIdentityState()
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "converting Identity State", err)
				return
			}

			if err := result.Identity.Set(ctx, *tfTypeIdentity); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "setting Identity Data", err)
				return
			}

			tfTypeResourceState, err := rd.TfTypeResourceState()
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "converting Resource State", err)
				return
			}

			if err := result.Resource.Set(ctx, *tfTypeResourceState); err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "setting Resource Data", err)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}

func (r *NetworkSecurityGroupListResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Defaults(request, response)
}

func (r *NetworkSecurityGroupListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, response *list.RawV5SchemaResponse) {
	res := resourceNetworkSecurityGroup()
	response.ProtoV5Schema = res.ProtoSchema(ctx)()
	response.ProtoV5IdentitySchema = res.ProtoIdentitySchema(ctx)()
}
