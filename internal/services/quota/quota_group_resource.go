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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotalimits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotassubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type QuotaRequestModel struct {
	ResourceName         string `tfschema:"resource_name"`
	Location             string `tfschema:"location"`
	ResourceProviderName string `tfschema:"resource_provider_name"`
	Limit                int64  `tfschema:"limit"`
	Comment              string `tfschema:"comment"`
	AvailableLimit       int64  `tfschema:"available_limit"`
}

type QuotaGroupModel struct {
	Name                      string              `tfschema:"name"`
	ManagementGroupId         string              `tfschema:"management_group_id"`
	DisplayName               string              `tfschema:"display_name"`
	AssociatedSubscriptionIds []string            `tfschema:"associated_subscription_ids"`
	QuotaRequests             []QuotaRequestModel `tfschema:"quota_request"`
}

var (
	_ sdk.Resource             = QuotaGroupResource{}
	_ sdk.ResourceWithUpdate   = QuotaGroupResource{}
	_ sdk.ResourceWithIdentity = QuotaGroupResource{}
)

type QuotaGroupResource struct{}

func (r QuotaGroupResource) ModelObject() interface{} {
	return &QuotaGroupModel{}
}

func (r QuotaGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return groupquotas.ValidateGroupQuotaID
}

func (r QuotaGroupResource) ResourceType() string {
	return "azurerm_quota_group"
}

func (r QuotaGroupResource) Identity() resourceids.ResourceId {
	return &groupquotas.GroupQuotaId{}
}

func (r QuotaGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
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
			ForceNew:     true,
			ValidateFunc: commonids.ValidateManagementGroupID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"associated_subscription_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"quota_request": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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

					"limit": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"comment": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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

func (r QuotaGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r QuotaGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			groupQuotasClient := metadata.Client.Quota.GroupQuotasClient
			subscriptionsClient := metadata.Client.Quota.GroupQuotasSubscriptionsClient
			limitsClient := metadata.Client.Quota.GroupQuotaLimitsClient

			var model QuotaGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			mgID, err := commonids.ParseManagementGroupID(model.ManagementGroupId)
			if err != nil {
				return fmt.Errorf("parsing management_group_id %q: %+v", model.ManagementGroupId, err)
			}

			id := groupquotas.NewGroupQuotaID(mgID.GroupId, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := groupQuotasClient.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := groupquotas.GroupQuotasEntity{
				Properties: &groupquotas.GroupQuotasEntityBase{},
			}
			if model.DisplayName != "" {
				payload.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if err := groupQuotasClient.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			for _, subID := range model.AssociatedSubscriptionIds {
				subscriptionID := groupquotassubscriptions.NewSubscriptionID(mgID.GroupId, model.Name, subID)
				if err := subscriptionsClient.GroupQuotaSubscriptionsCreateOrUpdateThenPoll(ctx, subscriptionID); err != nil {
					return fmt.Errorf("associating subscription %q with %s: %+v", subID, id, err)
				}
			}

			if err := applyQuotaRequests(ctx, limitsClient, id, model.QuotaRequests); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r QuotaGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			groupQuotasClient := metadata.Client.Quota.GroupQuotasClient
			subscriptionsClient := metadata.Client.Quota.GroupQuotasSubscriptionsClient
			limitsClient := metadata.Client.Quota.GroupQuotaLimitsClient

			id, err := groupquotas.ParseGroupQuotaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := groupQuotasClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			mgID := commonids.NewManagementGroupID(id.ManagementGroupId)

			subscriptionGroupQuotaID := groupquotassubscriptions.NewGroupQuotaID(id.ManagementGroupId, id.GroupQuotaName)

			subListResp, err := subscriptionsClient.GroupQuotaSubscriptionsListComplete(ctx, subscriptionGroupQuotaID)
			if err != nil {
				return fmt.Errorf("listing subscriptions for %s: %+v", *id, err)
			}

			state := QuotaGroupModel{
				Name:              id.GroupQuotaName,
				ManagementGroupId: mgID.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DisplayName = pointer.From(props.DisplayName)
				}
			}

			subscriptionIDs := make([]string, 0)
			for _, sub := range subListResp.Items {
				if sub.Properties != nil && sub.Properties.SubscriptionId != nil {
					subscriptionIDs = append(subscriptionIDs, *sub.Properties.SubscriptionId)
				}
			}
			state.AssociatedSubscriptionIds = subscriptionIDs

			quotaRequests, err := readQuotaRequests(ctx, limitsClient, id, metadata.ResourceData)
			if err != nil {
				return err
			}
			state.QuotaRequests = quotaRequests

			return metadata.Encode(&state)
		},
	}
}

func (r QuotaGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			groupQuotasClient := metadata.Client.Quota.GroupQuotasClient
			subscriptionsClient := metadata.Client.Quota.GroupQuotasSubscriptionsClient
			limitsClient := metadata.Client.Quota.GroupQuotaLimitsClient

			id, err := groupquotas.ParseGroupQuotaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model QuotaGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("display_name") {
				patch := groupquotas.GroupQuotasEntityPatch{
					Properties: &groupquotas.GroupQuotasEntityBasePatch{},
				}
				if model.DisplayName != "" {
					patch.Properties.DisplayName = pointer.To(model.DisplayName)
				} else {
					patch.Properties.DisplayName = pointer.To("")
				}

				if err := groupQuotasClient.UpdateThenPoll(ctx, *id, patch); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			if metadata.ResourceData.HasChange("associated_subscription_ids") {
				old, new := metadata.ResourceData.GetChange("associated_subscription_ids")
				oldSet := old.(*pluginsdk.Set)
				newSet := new.(*pluginsdk.Set)

				for _, v := range oldSet.Difference(newSet).List() {
					subID := v.(string)
					subscriptionID := groupquotassubscriptions.NewSubscriptionID(id.ManagementGroupId, id.GroupQuotaName, subID)
					if err := subscriptionsClient.GroupQuotaSubscriptionsDeleteThenPoll(ctx, subscriptionID); err != nil {
						return fmt.Errorf("disassociating subscription %q from %s: %+v", subID, *id, err)
					}
				}

				for _, v := range newSet.Difference(oldSet).List() {
					subID := v.(string)
					subscriptionID := groupquotassubscriptions.NewSubscriptionID(id.ManagementGroupId, id.GroupQuotaName, subID)
					if err := subscriptionsClient.GroupQuotaSubscriptionsCreateOrUpdateThenPoll(ctx, subscriptionID); err != nil {
						return fmt.Errorf("associating subscription %q with %s: %+v", subID, *id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("quota_request") {
				old, newVal := metadata.ResourceData.GetChange("quota_request")

				// Build a map of new (provider, location, resourceName) tuples so we can
				// detect which resources were present before but are no longer in config.
				type reqKey struct{ ResourceProvider, Location, ResourceName string }
				newReqSet := make(map[reqKey]bool)
				for _, raw := range newVal.(*pluginsdk.Set).List() {
					item := raw.(map[string]interface{})
					newReqSet[reqKey{
						ResourceProvider: item["resource_provider_name"].(string),
						Location:         location.Normalize(item["location"].(string)),
						ResourceName:     item["resource_name"].(string),
					}] = true
				}

				// For every resource_name that was in the old config but is absent from the
				// new config, send a zero-limit PATCH so Azure returns that quota to the pool.
				zeroRequests := make([]QuotaRequestModel, 0)
				for _, raw := range old.(*pluginsdk.Set).List() {
					item := raw.(map[string]interface{})
					rp := item["resource_provider_name"].(string)
					loc := location.Normalize(item["location"].(string))
					rn := item["resource_name"].(string)
					if !newReqSet[reqKey{ResourceProvider: rp, Location: loc, ResourceName: rn}] {
						zeroRequests = append(zeroRequests, QuotaRequestModel{
							ResourceName:         rn,
							Location:             loc,
							ResourceProviderName: rp,
							Limit:                0,
						})
					}
				}

				// Combine zeros (for removed resources) with the new desired state.
				allRequests := append(zeroRequests, model.QuotaRequests...)
				if err := applyQuotaRequests(ctx, limitsClient, *id, allRequests); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func (r QuotaGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			groupQuotasClient := metadata.Client.Quota.GroupQuotasClient
			subscriptionsClient := metadata.Client.Quota.GroupQuotasSubscriptionsClient

			id, err := groupquotas.ParseGroupQuotaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			subscriptionGroupQuotaID := groupquotassubscriptions.NewGroupQuotaID(id.ManagementGroupId, id.GroupQuotaName)

			subListResp, err := subscriptionsClient.GroupQuotaSubscriptionsListComplete(ctx, subscriptionGroupQuotaID)
			if err != nil {
				return fmt.Errorf("listing subscriptions for %s: %+v", *id, err)
			}

			for _, sub := range subListResp.Items {
				if sub.Properties == nil || sub.Properties.SubscriptionId == nil {
					continue
				}
				subscriptionID := groupquotassubscriptions.NewSubscriptionID(id.ManagementGroupId, id.GroupQuotaName, *sub.Properties.SubscriptionId)
				if err := subscriptionsClient.GroupQuotaSubscriptionsDeleteThenPoll(ctx, subscriptionID); err != nil {
					return fmt.Errorf("disassociating subscription %q from %s: %+v", *sub.Properties.SubscriptionId, *id, err)
				}
			}

			if err := groupQuotasClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

// applyQuotaRequests groups quota_request blocks by (resourceProvider, location) and submits
// a single RequestUpdate call per group, as required by the API.
func applyQuotaRequests(ctx context.Context, client *groupquotalimits.GroupQuotaLimitsClient, id groupquotas.GroupQuotaId, requests []QuotaRequestModel) error {
	type scopeKey struct {
		ResourceProvider string
		Location         string
	}

	grouped := make(map[scopeKey][]QuotaRequestModel)
	for _, req := range requests {
		key := scopeKey{
			ResourceProvider: req.ResourceProviderName,
			Location:         req.Location,
		}
		grouped[key] = append(grouped[key], req)
	}

	for key, reqs := range grouped {
		limitID := groupquotalimits.NewGroupQuotaLimitID(id.ManagementGroupId, id.GroupQuotaName, key.ResourceProvider, key.Location)

		items := make([]groupquotalimits.GroupQuotaLimit, 0, len(reqs))
		for _, req := range reqs {
			item := groupquotalimits.GroupQuotaLimit{
				Properties: &groupquotalimits.GroupQuotaDetails{
					ResourceName: pointer.To(req.ResourceName),
					Limit:        pointer.To(req.Limit),
				},
			}
			if req.Comment != "" {
				item.Properties.Comment = pointer.To(req.Comment)
			}
			items = append(items, item)
		}

		payload := groupquotalimits.GroupQuotaLimitList{
			Properties: &groupquotalimits.GroupQuotaLimitListProperties{
				Value: &items,
			},
		}

		if err := client.RequestUpdateThenPoll(ctx, limitID, payload); err != nil {
			return fmt.Errorf("submitting quota limit requests for %s (resource provider %q, location %q): %+v", id, key.ResourceProvider, key.Location, err)
		}
	}

	return nil
}

// readQuotaRequests reads back quota limit state for all (resourceProvider, location) scopes
// present in the current Terraform config, matching by resource_name.
func readQuotaRequests(ctx context.Context, client *groupquotalimits.GroupQuotaLimitsClient, id *groupquotas.GroupQuotaId, d *pluginsdk.ResourceData) ([]QuotaRequestModel, error) {
	type scopeKey struct {
		ResourceProvider string
		Location         string
	}

	// Build the set of scopes from config so we know what to read back.
	configScopes := make(map[scopeKey]map[string]QuotaRequestModel) // scope -> resourceName -> config
	if v, ok := d.GetOk("quota_request"); ok {
		for _, raw := range v.(*pluginsdk.Set).List() {
			item := raw.(map[string]interface{})
			rn := item["resource_name"].(string)
			loc := item["location"].(string)
			rp := item["resource_provider_name"].(string)
			key := scopeKey{ResourceProvider: rp, Location: loc}
			if configScopes[key] == nil {
				configScopes[key] = make(map[string]QuotaRequestModel)
			}
			configScopes[key][rn] = QuotaRequestModel{
				ResourceName:         rn,
				Location:             loc,
				ResourceProviderName: rp,
				Limit:                int64(item["limit"].(int)),
				Comment:              item["comment"].(string),
			}
		}
	}

	result := make([]QuotaRequestModel, 0)
	for key, configItems := range configScopes {
		limitID := groupquotalimits.NewGroupQuotaLimitID(id.ManagementGroupId, id.GroupQuotaName, key.ResourceProvider, key.Location)

		limits, httpResp, err := listAllGroupQuotaLimits(ctx, client, limitID)
		if err != nil {
			if response.WasNotFound(httpResp) {
				continue
			}
			return nil, fmt.Errorf("reading quota limits for %s (resource provider %q, location %q): %+v", *id, key.ResourceProvider, key.Location, err)
		}

		for _, limit := range limits {
			if limit.Properties == nil || limit.Properties.ResourceName == nil {
				continue
			}
			rn := *limit.Properties.ResourceName
			if _, wanted := configItems[rn]; !wanted {
				continue
			}
			qr := QuotaRequestModel{
				ResourceName:         rn,
				Location:             key.Location,
				ResourceProviderName: key.ResourceProvider,
				Limit:                pointer.From(limit.Properties.Limit),
				Comment:              pointer.From(limit.Properties.Comment),
				AvailableLimit:       pointer.From(limit.Properties.AvailableLimit),
			}
			result = append(result, qr)
		}
	}

	return result, nil
}

func (r QuotaGroupResource) flatten(metadata sdk.ResourceMetaData, id *groupquotas.GroupQuotaId, model *groupquotas.GroupQuotasEntity, subscriptionIDs []string) error {
	state := QuotaGroupModel{
		Name:              id.GroupQuotaName,
		ManagementGroupId: commonids.NewManagementGroupID(id.ManagementGroupId).ID(),
	}

	if model != nil {
		if props := model.Properties; props != nil {
			state.DisplayName = pointer.From(props.DisplayName)
		}
	}

	state.AssociatedSubscriptionIds = subscriptionIDs

	return metadata.Encode(&state)
}
