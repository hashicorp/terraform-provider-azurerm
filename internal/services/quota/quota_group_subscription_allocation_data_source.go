// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quota

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/subscriptionquotaallocation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type QuotaGroupSubscriptionAllocationDataSourceModel struct {
	QuotaGroupId         string            `tfschema:"quota_group_id"`
	SubscriptionId       string            `tfschema:"subscription_id"`
	Location             string            `tfschema:"location"`
	ResourceProviderName string            `tfschema:"resource_provider_name"`
	Allocations          []AllocationModel `tfschema:"allocation"`
}

var _ sdk.DataSource = QuotaGroupSubscriptionAllocationDataSource{}

type QuotaGroupSubscriptionAllocationDataSource struct{}

func (d QuotaGroupSubscriptionAllocationDataSource) ModelObject() interface{} {
	return &QuotaGroupSubscriptionAllocationDataSourceModel{}
}

func (d QuotaGroupSubscriptionAllocationDataSource) ResourceType() string {
	return "azurerm_quota_group_subscription_allocation"
}

func (d QuotaGroupSubscriptionAllocationDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"quota_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: groupquotas.ValidateGroupQuotaID,
		},

		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			StateFunc:    location.StateFunc,
		},

		"resource_provider_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "Microsoft.Compute",
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (d QuotaGroupSubscriptionAllocationDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"allocation": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"limit": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"shareable_quota": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d QuotaGroupSubscriptionAllocationDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Quota.SubscriptionQuotaAllocationClient

			var model QuotaGroupSubscriptionAllocationDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			groupID, err := groupquotas.ParseGroupQuotaID(model.QuotaGroupId)
			if err != nil {
				return fmt.Errorf("parsing quota_group_id %q: %+v", model.QuotaGroupId, err)
			}

			subID, err := commonids.ParseSubscriptionID(model.SubscriptionId)
			if err != nil {
				return fmt.Errorf("parsing subscription_id %q: %+v", model.SubscriptionId, err)
			}

			id := subscriptionquotaallocation.NewQuotaAllocationID(
				groupID.ManagementGroupId,
				subID.SubscriptionId,
				groupID.GroupQuotaName,
				model.ResourceProviderName,
				location.Normalize(model.Location),
			)

			allocs, httpResp, err := listAllSubscriptionAllocations(ctx, client, id)
			if err != nil {
				if response.WasNotFound(httpResp) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := QuotaGroupSubscriptionAllocationDataSourceModel{
				QuotaGroupId:         groupID.ID(),
				SubscriptionId:       commonids.NewSubscriptionID(subID.SubscriptionId).ID(),
				Location:             id.QuotaAllocationName,
				ResourceProviderName: id.ResourceProviderName,
			}

			allocations := make([]AllocationModel, 0)
			for _, item := range allocs {
				if item.Properties == nil || item.Properties.ResourceName == nil {
					continue
				}
				allocations = append(allocations, AllocationModel{
					ResourceName:   pointer.From(item.Properties.ResourceName),
					Limit:          pointer.From(item.Properties.Limit),
					ShareableQuota: pointer.From(item.Properties.ShareableQuota),
				})
			}
			state.Allocations = allocations

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
