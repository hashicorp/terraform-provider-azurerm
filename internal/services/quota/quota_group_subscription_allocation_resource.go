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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/subscriptionquotaallocation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AllocationModel struct {
	ResourceName   string `tfschema:"resource_name"`
	Limit          int64  `tfschema:"limit"`
	ShareableQuota int64  `tfschema:"shareable_quota"`
}

type QuotaGroupSubscriptionAllocationModel struct {
	QuotaGroupId         string            `tfschema:"quota_group_id"`
	SubscriptionId       string            `tfschema:"subscription_id"`
	Location             string            `tfschema:"location"`
	ResourceProviderName string            `tfschema:"resource_provider_name"`
	Allocations          []AllocationModel `tfschema:"allocation"`
}

var (
	_ sdk.Resource             = QuotaGroupSubscriptionAllocationResource{}
	_ sdk.ResourceWithUpdate   = QuotaGroupSubscriptionAllocationResource{}
	_ sdk.ResourceWithIdentity = QuotaGroupSubscriptionAllocationResource{}
)

type QuotaGroupSubscriptionAllocationResource struct{}

func (r QuotaGroupSubscriptionAllocationResource) ModelObject() interface{} {
	return &QuotaGroupSubscriptionAllocationModel{}
}

func (r QuotaGroupSubscriptionAllocationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return subscriptionquotaallocation.ValidateQuotaAllocationID
}

func (r QuotaGroupSubscriptionAllocationResource) ResourceType() string {
	return "azurerm_quota_group_subscription_allocation"
}

func (r QuotaGroupSubscriptionAllocationResource) Identity() resourceids.ResourceId {
	return &subscriptionquotaallocation.QuotaAllocationId{}
}

func (r QuotaGroupSubscriptionAllocationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"quota_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: groupquotas.ValidateGroupQuotaID,
		},

		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			StateFunc:    location.StateFunc,
		},

		"resource_provider_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "Microsoft.Compute",
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"allocation": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"limit": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
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

func (r QuotaGroupSubscriptionAllocationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r QuotaGroupSubscriptionAllocationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Quota.SubscriptionQuotaAllocationClient

			var model QuotaGroupSubscriptionAllocationModel
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GroupQuotaSubscriptionAllocationList(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.Value != nil {
					for _, item := range *existing.Model.Properties.Value {
						if item.Properties != nil && pointer.From(item.Properties.Limit) > 0 {
							return metadata.ResourceRequiresImport(r.ResourceType(), id)
						}
					}
				}
			}

			payload := buildAllocationPayload(model.Allocations)
			if err := client.GroupQuotaSubscriptionAllocationRequestUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r QuotaGroupSubscriptionAllocationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Quota.SubscriptionQuotaAllocationClient

			id, err := subscriptionquotaallocation.ParseQuotaAllocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GroupQuotaSubscriptionAllocationList(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			groupID := groupquotas.NewGroupQuotaID(id.ManagementGroupId, id.GroupQuotaName)
			subID := commonids.NewSubscriptionID(id.SubscriptionId)

			state := QuotaGroupSubscriptionAllocationModel{
				QuotaGroupId:         groupID.ID(),
				SubscriptionId:       subID.ID(),
				Location:             id.QuotaAllocationName,
				ResourceProviderName: id.ResourceProviderName,
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Value != nil {
				allocations := make([]AllocationModel, 0)
				for _, item := range *resp.Model.Properties.Value {
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
			}

			return metadata.Encode(&state)
		},
	}
}

func (r QuotaGroupSubscriptionAllocationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Quota.SubscriptionQuotaAllocationClient

			id, err := subscriptionquotaallocation.ParseQuotaAllocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model QuotaGroupSubscriptionAllocationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := buildAllocationPayload(model.Allocations)
			if err := client.GroupQuotaSubscriptionAllocationRequestUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r QuotaGroupSubscriptionAllocationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Quota.SubscriptionQuotaAllocationClient

			id, err := subscriptionquotaallocation.ParseQuotaAllocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Read current allocations so we can zero them out to return quota to the pool.
			resp, err := client.GroupQuotaSubscriptionAllocationList(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("reading %s for deletion: %+v", *id, err)
			}

			if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Value == nil {
				return nil
			}

			// Zero out all limits to return quota to the group pool.
			zeroed := make([]AllocationModel, 0, len(*resp.Model.Properties.Value))
			for _, item := range *resp.Model.Properties.Value {
				if item.Properties == nil || item.Properties.ResourceName == nil {
					continue
				}
				zeroed = append(zeroed, AllocationModel{
					ResourceName: pointer.From(item.Properties.ResourceName),
					Limit:        0,
				})
			}

			if len(zeroed) == 0 {
				return nil
			}

			payload := buildAllocationPayload(zeroed)
			if err := client.GroupQuotaSubscriptionAllocationRequestUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("zeroing allocations for %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func buildAllocationPayload(allocations []AllocationModel) subscriptionquotaallocation.SubscriptionQuotaAllocationsList {
	items := make([]subscriptionquotaallocation.SubscriptionQuotaAllocations, 0, len(allocations))
	for _, a := range allocations {
		items = append(items, subscriptionquotaallocation.SubscriptionQuotaAllocations{
			Properties: &subscriptionquotaallocation.SubscriptionQuotaDetails{
				ResourceName: pointer.To(a.ResourceName),
				Limit:        pointer.To(a.Limit),
				// ShareableQuota is computed by the API (limit minus current usage);
				// it is not accepted as a request field.
			},
		})
	}
	return subscriptionquotaallocation.SubscriptionQuotaAllocationsList{
		Properties: &subscriptionquotaallocation.SubscriptionQuotaAllocationsListProperties{
			Value: &items,
		},
	}
}
