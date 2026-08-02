// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quota

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotalimits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotassubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type QuotaGroupDataSourceModel struct {
	Name                      string              `tfschema:"name"`
	ManagementGroupId         string              `tfschema:"management_group_id"`
	Location                  string              `tfschema:"location"`
	ResourceProviderName      string              `tfschema:"resource_provider_name"`
	DisplayName               string              `tfschema:"display_name"`
	AssociatedSubscriptionIds []string            `tfschema:"associated_subscription_ids"`
	QuotaRequests             []QuotaRequestModel `tfschema:"quota_request"`
}

var _ sdk.DataSource = QuotaGroupDataSource{}

type QuotaGroupDataSource struct{}

func (d QuotaGroupDataSource) ModelObject() interface{} {
	return &QuotaGroupDataSourceModel{}
}

func (d QuotaGroupDataSource) ResourceType() string {
	return "azurerm_quota_group"
}

func (d QuotaGroupDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringMatch(
					regexp.MustCompile(`^[a-z][a-z0-9]*$`),
					"must start with a lowercase letter and contain only lowercase letters and digits",
				),
				validation.StringLenBetween(3, 63),
			),
		},

		"management_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateManagementGroupID,
		},

		"location": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
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

func (d QuotaGroupDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"associated_subscription_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"quota_request": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"location": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"resource_provider_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"limit": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"comment": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"available_limit": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func (d QuotaGroupDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			groupQuotasClient := metadata.Client.Quota.GroupQuotasClient
			subscriptionsClient := metadata.Client.Quota.GroupQuotasSubscriptionsClient
			limitsClient := metadata.Client.Quota.GroupQuotaLimitsClient

			var model QuotaGroupDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			mgID, err := commonids.ParseManagementGroupID(model.ManagementGroupId)
			if err != nil {
				return fmt.Errorf("parsing management_group_id %q: %+v", model.ManagementGroupId, err)
			}

			id := groupquotas.NewGroupQuotaID(mgID.GroupId, model.Name)

			resp, err := groupQuotasClient.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := QuotaGroupDataSourceModel{
				Name:                 id.GroupQuotaName,
				ManagementGroupId:    commonids.NewManagementGroupID(id.ManagementGroupId).ID(),
				Location:             model.Location,
				ResourceProviderName: model.ResourceProviderName,
			}

			if resp.Model != nil && resp.Model.Properties != nil {
				state.DisplayName = pointer.From(resp.Model.Properties.DisplayName)
			}

			subscriptionGroupQuotaID := groupquotassubscriptions.NewGroupQuotaID(id.ManagementGroupId, id.GroupQuotaName)
			subListResp, err := subscriptionsClient.GroupQuotaSubscriptionsListComplete(ctx, subscriptionGroupQuotaID)
			if err != nil {
				return fmt.Errorf("listing subscriptions for %s: %+v", id, err)
			}
			subIDs := make([]string, 0)
			for _, sub := range subListResp.Items {
				if sub.Properties != nil && sub.Properties.SubscriptionId != nil {
					subIDs = append(subIDs, *sub.Properties.SubscriptionId)
				}
			}
			state.AssociatedSubscriptionIds = subIDs

			// If location is provided, fetch all quota_request items for that (resource_provider, location)
			// scope, following NextLink pagination.
			if loc := location.Normalize(model.Location); loc != "" {
				limitID := groupquotalimits.NewGroupQuotaLimitID(id.ManagementGroupId, id.GroupQuotaName, model.ResourceProviderName, loc)
				limits, httpResp, err := listAllGroupQuotaLimits(ctx, limitsClient, limitID)
				if err != nil {
					if !response.WasNotFound(httpResp) {
						return fmt.Errorf("reading quota limits for %s (resource provider %q, location %q): %+v", id, model.ResourceProviderName, loc, err)
					}
				} else {
					quotaRequests := make([]QuotaRequestModel, 0)
					for _, limit := range limits {
						if limit.Properties == nil || limit.Properties.ResourceName == nil {
							continue
						}
						quotaRequests = append(quotaRequests, QuotaRequestModel{
							ResourceName:         pointer.From(limit.Properties.ResourceName),
							Location:             loc,
							ResourceProviderName: model.ResourceProviderName,
							Limit:                pointer.From(limit.Properties.Limit),
							Comment:              pointer.From(limit.Properties.Comment),
							AvailableLimit:       pointer.From(limit.Properties.AvailableLimit),
						})
					}
					state.QuotaRequests = quotaRequests
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
