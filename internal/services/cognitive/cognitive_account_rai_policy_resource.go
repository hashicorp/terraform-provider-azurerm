// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithUpdate        = &CognitiveAccountRaiPolicyResource{}
	_ sdk.ResourceWithCustomizeDiff = &CognitiveAccountRaiPolicyResource{}
)

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

var severityThresholdNotApplicableFilterNames = []string{
	"Jailbreak",
	"Indirect Attack",
	"Protected Material Text",
	"Protected Material Code",
}

func (r CognitiveAccountRaiPolicyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			rawFilters, ok := metadata.ResourceDiff.GetOk("content_filter")
			if !ok {
				return nil
			}

			filters, ok := rawFilters.([]interface{})
			if !ok {
				return nil
			}

			for i, rawFilter := range filters {
				filter, ok := rawFilter.(map[string]interface{})
				if !ok {
					continue
				}

				name, _ := filter["name"].(string)
				severityThreshold, _ := filter["severity_threshold"].(string)

				if severityThreshold == "" {
					continue
				}

				for _, notApplicable := range severityThresholdNotApplicableFilterNames {
					if name == notApplicable {
						return fmt.Errorf("`severity_threshold` is not applicable for `content_filter[%d]` with name %q", i, name)
					}
				}
			}

			return nil
		},
	}
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
					"block_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"filter_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"source": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForRaiPolicyContentSource(), false),
					},
					"severity_threshold": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForContentLevel(), false),
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

			id := raipolicies.NewAccountRaiPolicyID(subscriptionId, cognitiveAccountId.ResourceGroupName, cognitiveAccountId.AccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

			return nil
		},
	}
}

func (r CognitiveAccountRaiPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient

			id, err := raipolicies.ParseAccountRaiPolicyID(metadata.ResourceData.Id())
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

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountRaiPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient

			id, err := raipolicies.ParseAccountRaiPolicyID(metadata.ResourceData.Id())
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

			id, err := raipolicies.ParseAccountRaiPolicyID(metadata.ResourceData.Id())
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
	return raipolicies.ValidateAccountRaiPolicyID
}

func expandRaiPolicyContentFilters(filters []AccountRaiPolicyContentFilter) *[]raipolicies.RaiPolicyContentFilter {
	if filters == nil {
		return nil
	}

	contentFilters := make([]raipolicies.RaiPolicyContentFilter, 0, len(filters))
	for _, filter := range filters {
		f := raipolicies.RaiPolicyContentFilter{
			Name:     pointer.To(filter.Name),
			Enabled:  pointer.To(filter.FilterEnabled),
			Blocking: pointer.To(filter.BlockEnabled),
			Source:   pointer.To(raipolicies.RaiPolicyContentSource(filter.Source)),
		}

		if filter.SeverityThreshold != "" {
			f.SeverityThreshold = pointer.To(raipolicies.ContentLevel(filter.SeverityThreshold))
		}

		contentFilters = append(contentFilters, f)
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
