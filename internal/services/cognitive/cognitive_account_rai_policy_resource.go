// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cognitive_account_rai_policy -properties "name,cognitive_account_id" -service-package-name cognitive -known-values "subscription_id:data.Subscriptions.Primary"

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = &CognitiveAccountRaiPolicyResource{}
var _ sdk.ResourceWithIdentity = CognitiveAccountRaiPolicyResource{}

type CognitiveAccountRaiPolicyResource struct{}

type AccountRaiPolicyContentFilter struct {
	Name              string `tfschema:"name"`
	FilterEnabled     bool   `tfschema:"filter_enabled"`
	BlockEnabled      bool   `tfschema:"block_enabled"`
	SeverityThreshold string `tfschema:"severity_threshold"`
	Source            string `tfschema:"source"`
}

type AccountRaiPolicyCustomBlock struct {
	Id           string `tfschema:"rai_blocklist_id"`
	BlockEnabled bool   `tfschema:"block_enabled"`
	Source       string `tfschema:"source"`
}

type AccountRaiPolicyResourceModel struct {
	Name           string                          `tfschema:"name"`
	AccountId      string                          `tfschema:"cognitive_account_id"`
	BasePolicyName string                          `tfschema:"base_policy_name"`
	ContentFilter  []AccountRaiPolicyContentFilter `tfschema:"content_filter"`
	Mode           string                          `tfschema:"mode"`
	Tags           map[string]string               `tfschema:"tags"`
}

func (r CognitiveAccountRaiPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: raipolicies.ValidateAccountID,
		},

		"base_policy_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"content_filter": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"filter_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"block_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"severity_threshold": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForContentLevel(), false),
					},
					"source": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForRaiPolicyContentSource(), false),
					},
				},
			},
		},

		"mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForRaiPolicyMode(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r CognitiveAccountRaiPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountRaiPolicyResource) ModelObject() interface{} {
	return &AccountRaiPolicyResourceModel{}
}

func (r CognitiveAccountRaiPolicyResource) ResourceType() string {
	return "azurerm_cognitive_account_rai_policy"
}

func (r CognitiveAccountRaiPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AccountRaiPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			cognitiveAccountId, err := raipolicies.ParseAccountID(model.AccountId)
			if err != nil {
				return err
			}

			id := raipolicies.NewRaiPolicyID(subscriptionId, cognitiveAccountId.ResourceGroupName, cognitiveAccountId.AccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			locks.ByID(cognitiveAccountId.ID())
			defer locks.UnlockByID(cognitiveAccountId.ID())

			raiPolicy := raipolicies.RaiPolicy{
				Name: pointer.To(model.Name),
				Properties: &raipolicies.RaiPolicyProperties{
					BasePolicyName: pointer.To(model.BasePolicyName),
					ContentFilters: expandRaiPolicyContentFilters(model.ContentFilter),
				},
				Tags: pointer.To(model.Tags),
			}

			if model.Mode != "" {
				raiPolicy.Properties.Mode = pointer.To(raipolicies.RaiPolicyMode(model.Mode))
			}

			if _, err := client.CreateOrUpdate(ctx, id, raiPolicy); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r CognitiveAccountRaiPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient

			id, err := raipolicies.ParseRaiPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			cognitiveAccountId := raipolicies.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AccountRaiPolicyResourceModel{
				Name:      id.RaiPolicyName,
				AccountId: cognitiveAccountId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.BasePolicyName = pointer.From(props.BasePolicyName)
					state.ContentFilter = flattenRaiPolicyContentFilters(props.ContentFilters)
					state.Mode = string(pointer.From(props.Mode))
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountRaiPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient

			id, err := raipolicies.ParseRaiPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AccountRaiPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			cognitiveAccountId := raipolicies.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(cognitiveAccountId.ID())
			defer locks.UnlockByID(cognitiveAccountId.ID())

			payload := existing.Model

			if metadata.ResourceData.HasChange("content_filter") {
				payload.Properties.ContentFilters = expandRaiPolicyContentFilters(model.ContentFilter)
			}

			if metadata.ResourceData.HasChange("mode") {
				payload.Properties.Mode = pointer.To(raipolicies.RaiPolicyMode(model.Mode))
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountRaiPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient

			id, err := raipolicies.ParseRaiPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			cognitiveAccountId := raipolicies.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(cognitiveAccountId.ID())
			defer locks.UnlockByID(cognitiveAccountId.ID())

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountRaiPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return raipolicies.ValidateRaiPolicyID
}

func (r CognitiveAccountRaiPolicyResource) Identity() resourceids.ResourceId {
	return &raipolicies.RaiPolicyId{}
}

func expandRaiPolicyContentFilters(filters []AccountRaiPolicyContentFilter) *[]raipolicies.RaiPolicyContentFilter {
	if filters == nil {
		return nil
	}

	contentFilters := make([]raipolicies.RaiPolicyContentFilter, 0, len(filters))
	for _, filter := range filters {
		contentFilters = append(contentFilters, raipolicies.RaiPolicyContentFilter{
			Name:              pointer.To(filter.Name),
			Enabled:           pointer.To(filter.FilterEnabled),
			Blocking:          pointer.To(filter.BlockEnabled),
			SeverityThreshold: pointer.To(raipolicies.ContentLevel(filter.SeverityThreshold)),
			Source:            pointer.To(raipolicies.RaiPolicyContentSource(filter.Source)),
		})
	}
	return &contentFilters
}

func flattenRaiPolicyContentFilters(filters *[]raipolicies.RaiPolicyContentFilter) []AccountRaiPolicyContentFilter {
	contentFilters := make([]AccountRaiPolicyContentFilter, 0)
	if filters == nil {
		return contentFilters
	}

	for _, filter := range *filters {
		contentFilters = append(contentFilters, AccountRaiPolicyContentFilter{
			Name:              pointer.From(filter.Name),
			FilterEnabled:     pointer.From(filter.Enabled),
			BlockEnabled:      pointer.From(filter.Blocking),
			SeverityThreshold: string(pointer.From(filter.SeverityThreshold)),
			Source:            string(pointer.From(filter.Source)),
		})
	}
	return contentFilters
}
