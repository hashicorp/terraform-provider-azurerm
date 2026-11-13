// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package quota

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotassubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/subscriptionquotaallocation"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type QuotaGroupSubscriptionAllocationListResource struct{}

// The Quota API has no operation that enumerates allocations across providers and locations, so a
// list is scoped to a single (quota group, resource provider, location) and walks the group's
// associated subscriptions.
type QuotaGroupSubscriptionAllocationListModel struct {
	QuotaGroupId         types.String `tfsdk:"quota_group_id"`
	Location             types.String `tfsdk:"location"`
	ResourceProviderName types.String `tfsdk:"resource_provider_name"`
}

var _ sdk.FrameworkListWrappedResource = new(QuotaGroupSubscriptionAllocationListResource)

func (QuotaGroupSubscriptionAllocationListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = QuotaGroupSubscriptionAllocationResource{}.ResourceType()
}

func (QuotaGroupSubscriptionAllocationListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(QuotaGroupSubscriptionAllocationResource{})
}

func (QuotaGroupSubscriptionAllocationListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"quota_group_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: groupquotas.ValidateGroupQuotaID},
				},
			},

			"location": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: validation.StringIsNotEmpty},
				},
			},

			"resource_provider_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: validation.StringIsNotEmpty},
				},
			},
		},
	}
}

func (QuotaGroupSubscriptionAllocationListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	subscriptionsClient := metadata.Client.Quota.GroupQuotasSubscriptionsClient
	allocationClient := metadata.Client.Quota.SubscriptionQuotaAllocationClient

	var data QuotaGroupSubscriptionAllocationListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	r := QuotaGroupSubscriptionAllocationResource{}

	groupID, err := groupquotas.ParseGroupQuotaID(data.QuotaGroupId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing quota_group_id for `%s`", r.ResourceType()), err)
		return
	}

	resourceProviderName := "Microsoft.Compute"
	if !data.ResourceProviderName.IsNull() && !data.ResourceProviderName.IsUnknown() {
		resourceProviderName = data.ResourceProviderName.ValueString()
	}
	locationName := location.Normalize(data.Location.ValueString())

	subscriptions, err := subscriptionsClient.GroupQuotaSubscriptionsListComplete(ctx, groupquotassubscriptions.NewGroupQuotaID(groupID.ManagementGroupId, groupID.GroupQuotaName))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing subscriptions for `%s`", r.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, subscription := range subscriptions.Items {
			if subscription.Properties == nil || subscription.Properties.SubscriptionId == nil {
				continue
			}

			result := request.NewListResult(ctx)
			result.DisplayName = fmt.Sprintf("%s/%s", pointer.From(subscription.Properties.SubscriptionId), locationName)

			id := subscriptionquotaallocation.NewQuotaAllocationID(groupID.ManagementGroupId, pointer.From(subscription.Properties.SubscriptionId), groupID.GroupQuotaName, resourceProviderName, locationName)

			allocations, _, err := listAllSubscriptionAllocations(ctx, allocationClient, id)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("listing allocations for %s", id), err)
				return
			}

			// An allocation only exists when at least one resource has a non-zero limit, matching the check in Create.
			hasAllocation := false
			for _, item := range allocations {
				if item.Properties != nil && pointer.From(item.Properties.Limit) > 0 {
					hasAllocation = true
					break
				}
			}
			if !hasAllocation {
				continue
			}

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)

			if err := r.flatten(meta, &id, allocations); err != nil {
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
