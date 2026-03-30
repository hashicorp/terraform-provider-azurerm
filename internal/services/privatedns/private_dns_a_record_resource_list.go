// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package privatedns

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type (
	PrivateDnsARecordListResource struct{}
	PrivateDnsARecordListModel    struct {
		PrivateDnsZoneId types.String `tfsdk:"private_dns_zone_id"`
	}
)

var _ sdk.FrameworkListWrappedResource = new(PrivateDnsARecordListResource)

func (r PrivateDnsARecordListResource) ResourceFunc() *pluginsdk.Resource {
	return resourcePrivateDnsARecord()
}

func (r PrivateDnsARecordListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_private_dns_a_record"
}

func (r PrivateDnsARecordListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"private_dns_zone_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: privatezones.ValidatePrivateDnsZoneID,
					},
				},
			},
		},
	}
}

func (r PrivateDnsARecordListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.PrivateDns.RecordSetsClient
	var data PrivateDnsARecordListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	results := make([]privatedns.RecordSet, 0)

	if !data.PrivateDnsZoneId.IsNull() {
		privatezoneId, err := privatezones.ParsePrivateDnsZoneID(data.PrivateDnsZoneId.ValueString())

		privateDnsZoneId := privatedns.NewPrivateZoneID(privatezoneId.SubscriptionId, privatezoneId.ResourceGroupName, privatezoneId.PrivateDnsZoneName, privatedns.RecordTypeA)

		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing privatedns PrivateDnsZone ID for `%s`", "azurerm_private_dns_a_record"), err)
			return
		}

		resp, err := client.RecordSetsListByTypeComplete(ctx, privateDnsZoneId, privatedns.DefaultRecordSetsListByTypeOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_private_dns_a_record"), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, arecord := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(arecord.Name)

			rd := resourcePrivateDnsARecord().Data(&terraform.InstanceState{})

			id, err := privatedns.ParseRecordTypeID(pointer.From(arecord.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing PrivateDns ARecord ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourcePrivateDnsARecordFlatten(rd, id, &arecord); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_private_dns_a_record"), err)
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
