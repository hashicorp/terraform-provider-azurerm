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

type DnsTxtRecordListResource struct{}

type DnsTxtRecordListModel struct {
	DnsZoneId types.String `tfsdk:"dns_zone_id"`
}

var _ sdk.FrameworkListWrappedResource = new(DnsTxtRecordListResource)

func (DnsTxtRecordListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azurermDnsTxtRecordResourceName
}

func (DnsTxtRecordListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceDnsTxtRecord()
}

func (DnsTxtRecordListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
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

func (DnsTxtRecordListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Dns.RecordSets

	var data DnsTxtRecordListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	dnsZoneId, err := recordsets.ParseDnsZoneID(data.DnsZoneId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing parent ID for `%s`", azurermDnsTxtRecordResourceName), err)
		return
	}

	zoneId := recordsets.NewZoneID(dnsZoneId.SubscriptionId, dnsZoneId.ResourceGroupName, dnsZoneId.DnsZoneName, recordsets.RecordTypeTXT)

	resp, err := client.ListByTypeComplete(ctx, zoneId, recordsets.DefaultListByTypeOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azurermDnsTxtRecordResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceDnsTxtRecord().Data(&terraform.InstanceState{})

			id, err := recordsets.ParseRecordTypeIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing %s ID", azurermDnsTxtRecordResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceDnsTxtRecordFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azurermDnsTxtRecordResourceName), err)
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
