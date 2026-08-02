// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quota

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotassubscriptions"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type QuotaGroupListResource struct{}

type QuotaGroupListModel struct {
	ManagementGroupId types.String `tfsdk:"management_group_id"`
}

var _ sdk.FrameworkListWrappedResource = new(QuotaGroupListResource)

func (QuotaGroupListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = QuotaGroupResource{}.ResourceType()
}

func (QuotaGroupListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(QuotaGroupResource{})
}

func (QuotaGroupListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"management_group_id": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: commonids.ValidateManagementGroupID},
				},
			},
		},
	}
}

func (QuotaGroupListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Quota.GroupQuotasClient
	subscriptionsClient := metadata.Client.Quota.GroupQuotasSubscriptionsClient

	var data QuotaGroupListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	mgID, err := commonids.ParseManagementGroupID(data.ManagementGroupId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing management_group_id for `%s`", QuotaGroupResource{}.ResourceType()), err)
		return
	}

	resp, err := client.ListComplete(ctx, *mgID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", QuotaGroupResource{}.ResourceType()), err)
		return
	}

	r := QuotaGroupResource{}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, quotaGroup := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(quotaGroup.Name)

			id, err := groupquotas.ParseGroupQuotaIDInsensitively(pointer.From(quotaGroup.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Quota Group ID", err)
				return
			}

			// Fetch associated subscriptions for this quota group.
			subscriptionGroupQuotaID := groupquotassubscriptions.NewGroupQuotaID(id.ManagementGroupId, id.GroupQuotaName)
			subListResp, err := subscriptionsClient.GroupQuotaSubscriptionsListComplete(ctx, subscriptionGroupQuotaID)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("listing subscriptions for %s", id), err)
				return
			}

			subscriptionIDs := make([]string, 0)
			for _, sub := range subListResp.Items {
				if sub.Properties != nil && sub.Properties.SubscriptionId != nil {
					subscriptionIDs = append(subscriptionIDs, *sub.Properties.SubscriptionId)
				}
			}

			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)

			if err := r.flatten(meta, id, &quotaGroup, subscriptionIDs); err != nil {
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
