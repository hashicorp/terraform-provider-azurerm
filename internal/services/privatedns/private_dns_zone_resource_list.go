package privatedns

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateDnsZoneListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(PrivateDnsZoneListResource)

func (r PrivateDnsZoneListResource) ResourceFunc() *pluginsdk.Resource {
	return resourcePrivateDnsZone()
}

func (r PrivateDnsZoneListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = privateDnsZoneResourceName
}

func (r PrivateDnsZoneListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.PrivateDns.PrivateZonesClient
	recordSetsClient := metadata.Client.PrivateDns.RecordSetsClient

	ctx, cancel := context.WithTimeout(ctx, time.Minute*60)
	defer cancel()

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]privatezones.PrivateZone, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroup(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()), privatezones.DefaultListByResourceGroupOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", privateDnsZoneResourceName), err)
			return
		}

		results = pointer.From(resp.Model)
	default:
		resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID), privatezones.DefaultListOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", privateDnsZoneResourceName), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, privateZone := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(privateZone.Name)

			id, err := privatezones.ParsePrivateDnsZoneID(pointer.From(privateZone.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Private DNS Zone ID", err)
				return
			}

			rd := resourcePrivateDnsZone().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			recordId := privatedns.NewRecordTypeID(id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName, privatedns.RecordTypeSOA, "@")
			recordSetResp, err := recordSetsClient.RecordSetsGet(ctx, recordId)
			if err != nil {
				sdk.SetResponseErrorDiagnostic(stream, "reading DNS SOA record @: %v", err)
				return
			}

			if err := resourcePrivateDnsZoneFlatten(rd, id, &privateZone, recordSetResp.Model); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("encoding `%s` resource data", privateDnsZoneResourceName), err)
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
