package dns

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DnsCNameRecordListResource struct{}

type DnsCNameRecordListModel struct {
	DnsZoneId types.String `tfsdk:"dns_zone_id"`
}

var _ sdk.FrameworkListWrappedResource = new(DnsCNameRecordListResource)

func (DnsCNameRecordListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureDnsCNameRecordResourceName
}

func (DnsCNameRecordListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceDnsCNameRecord()
}

func (DnsCNameRecordListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dns_zone_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: recordsets.ValidateDnsZoneID},
				},
			},
		},
	}
}

func (DnsCNameRecordListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Dns.RecordSets

	var data DnsCNameRecordListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := recordsets.ParseDnsZoneID(data.DnsZoneId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing parent ID for `%s`", azureDnsCNameRecordResourceName), err)
		return
	}

	cnameType := string(recordsets.RecordTypeCNAME)
	resp, err := client.ListByDnsZoneCompleteMatchingPredicate(ctx, *parentID, recordsets.DefaultListByDnsZoneOperationOptions(), recordsets.RecordSetOperationPredicate{
		Type: &cnameType,
	})
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureDnsCNameRecordResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceDnsCNameRecord().Data(&terraform.InstanceState{})
			id, err := recordsets.ParseRecordTypeIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", azureDnsCNameRecordResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceDnsCNameRecordFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureDnsCNameRecordResourceName), err)
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
