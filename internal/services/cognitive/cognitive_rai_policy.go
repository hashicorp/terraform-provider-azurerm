package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/raipolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type cognitiveRaiPolicyModel struct {
	Name               string                        `tfschema:"name"`
	CognitiveAccountId string                        `tfschema:"cognitive_account_id"`
	BasePolicyName     string                        `tfschema:"base_policy_name"`
	Mode               string                        `tfschema:"mode"`
	Type               string                        `tfschema:"type"`
	ContentFilters     []RaiPolicyContentFilterModel `tfschema:"content_filters"`
	CustomBlocklists   []CustomBlocklistsModel       `tfschema:"custom_blocklists"`
}

type RaiPolicyContentFilterModel struct {
	Name              string `tfschema:"name"`
	Blocking          bool   `tfschema:"blocking"`
	Enabled           bool   `tfschema:"enabled"`
	SeverityThreshold string `tfschema:"severity_threshold"`
	Source            string `tfschema:"source"`
}

type CustomBlocklistsModel struct {
	Blocking      bool   `tfschema:"blocking"`
	BlocklistName string `tfschema:"blocklists_name"`
	Source        string `tfschema:"source"`
}

type CognitiveRaiPolicyResource struct{}

var _ sdk.Resource = CognitiveRaiPolicyResource{}

func (c CognitiveRaiPolicyResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// ForceNew:     true,
			ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
		},

		"base_policy_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Microsoft.Default",
			}, false),
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Default",
				"Asynchronous_filter",
				"Blocking",
				"Deferred",
			}, false),
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"SystemManaged",
				"UserManaged",
			}, false),
		},

		"content_filters": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			// ForceNew: true,
			// MinItems: 1,
			MaxItems: 8,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"hate",
							"sexual",
							"selfharm",
							"violence",
							"jailbreak",
							"protected_material_text",
							"protected_material_code",
							"profanity",
						}, false),
					},

					"blocking": {
						Type:     pluginsdk.TypeBool,
						Required: true,
						// ForceNew: true,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
						// ForceNew: true,
					},

					"severity_threshold": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						// ForceNew: true,
						Description: "Medium as default value",
						ValidateFunc: validation.StringInSlice([]string{
							"High",
							"Medium",
							"Low",
						}, false),
					},

					"source": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Prompt",
							"Completion",
						}, false),
					},
				},
			},
		},

		"custom_blocklists": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"blocklist_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// ForceNew:     true,
						ValidateFunc: validation.StringIsEmpty,
					},

					"blocking": {
						Type:     pluginsdk.TypeBool,
						Required: true,
						// ForceNew: true,
					},

					"source": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Prompt",
							"Completion",
						}, false),
					},
				},
			},
		},
	}
}

func (c CognitiveRaiPolicyResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (c CognitiveRaiPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model cognitiveRaiPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cognitive.RaiPoliciesClient
			accountId, err := raipolicies.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id := raipolicies.NewRaiPolicyID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(c.ResourceType(), id)
			}

			properties := &raipolicies.RaiPolicy{
				Properties: &raipolicies.RaiPolicyProperties{},
			}

			if model.BasePolicyName != "" {
				properties.Properties.BasePolicyName = &model.BasePolicyName
			}

			if model.Mode != "" {
				mode := raipolicies.RaiPolicyMode(model.Mode)
				properties.Properties.Mode = &mode
			}

			if model.Type != "" {
				rType := raipolicies.RaiPolicyType(model.Type)
				properties.Properties.Type = &rType
			}

			properties.Properties.ContentFilters = expandContentFiltersModel(model.ContentFilters)
			properties.Properties.CustomBlocklists = expandCustomBlocklistsModel(model.CustomBlocklists)

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (c CognitiveRaiPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model cognitiveRaiPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cognitive.RaiPoliciesClient
			accountId, err := raipolicies.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id, err := raipolicies.ParseRaiPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}

			properties := existing.Model

			if metadata.ResourceData.HasChange("base_policy_name") {
				properties.Properties.BasePolicyName = &model.BasePolicyName
			}

			if metadata.ResourceData.HasChange("mode") {
				mode := raipolicies.RaiPolicyMode(model.Mode)
				properties.Properties.Mode = &mode
			}

			if metadata.ResourceData.HasChange("type") {
				rType := raipolicies.RaiPolicyType(model.Type)
				properties.Properties.Type = &rType
			}

			if metadata.ResourceData.HasChange("content_filters") {
				return nil
			}

			if metadata.ResourceData.HasChange("custom_blocklists") {
				return nil
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (c CognitiveRaiPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient

			id, err := raipolicies.ParseRaiPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			accountId := cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (c CognitiveRaiPolicyResource) IDValidationFunc() func(interface{}, string) ([]string, []error) {
	return raipolicies.ValidateRaiPolicyID
}

func (c CognitiveRaiPolicyResource) ModelObject() interface{} {
	return &cognitiveRaiPolicyModel{}
}

func (c CognitiveRaiPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient

			id, err := raipolicies.ParseRaiPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := cognitiveRaiPolicyModel{
				Name:               id.RaiPolicyName,
				CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			if properties := model.Properties; properties != nil {
				if r := properties.BasePolicyName; r != nil {
					state.BasePolicyName = *r
				}
				if r := properties.Mode; r != nil {
					state.Mode = string(*r)
				}
				if r := properties.Type; r != nil {
					state.Type = string(*r)
				}

				state.ContentFilters = flattenContentFiltersModel(properties.ContentFilters)

				if customBlockLists := flattenCustomBlicklistsModel(properties.CustomBlocklists); customBlockLists != nil {
					state.CustomBlocklists = customBlockLists
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (c CognitiveRaiPolicyResource) ResourceType() string {
	return "azurerm_cognitive_rai_policy"
}

func expandContentFiltersModel(inputList []RaiPolicyContentFilterModel) *[]raipolicies.RaiPolicyContentFilter {
	contentFilters := make([]raipolicies.RaiPolicyContentFilter, 0)

	for _, contentFilter := range inputList {
		contentFilterRaw := raipolicies.RaiPolicyContentFilter{
			Blocking: &contentFilter.Blocking,
			Enabled:  &contentFilter.Enabled,
		}

		if contentFilter.Name != "" {
			contentFilterRaw.Name = &contentFilter.Name
		}

		if contentFilter.Source != "" {
			source := raipolicies.RaiPolicyContentSource(contentFilter.Source)
			contentFilterRaw.Source = &source
		}

		if contentFilter.SeverityThreshold != "" {
			sThreshold := raipolicies.ContentLevel(contentFilter.SeverityThreshold)
			contentFilterRaw.SeverityThreshold = &sThreshold
		}

		contentFilters = append(contentFilters, contentFilterRaw)
	}

	return &contentFilters
}

func expandCustomBlocklistsModel(inputList []CustomBlocklistsModel) *[]raipolicies.CustomBlocklistConfig {
	customBlocklists := make([]raipolicies.CustomBlocklistConfig, 0)

	for _, customBlocklist := range inputList {
		customBlocklistRaw := raipolicies.CustomBlocklistConfig{
			Blocking: &customBlocklist.Blocking,
		}

		if customBlocklist.BlocklistName != "" {
			customBlocklistRaw.BlocklistName = &customBlocklist.BlocklistName
		}

		if customBlocklist.Source != "" {
			source := raipolicies.RaiPolicyContentSource(customBlocklist.Source)
			customBlocklistRaw.Source = &source
		}

		customBlocklists = append(customBlocklists, customBlocklistRaw)
	}

	return &customBlocklists
}

func flattenContentFiltersModel(inputList *[]raipolicies.RaiPolicyContentFilter) []RaiPolicyContentFilterModel {
	contentFilters := make([]RaiPolicyContentFilterModel, 0)

	if inputList == nil {
		return contentFilters
	}

	for _, contentFilter := range *inputList {
		contentFilterRaw := RaiPolicyContentFilterModel{}

		if contentFilter.Blocking != nil {
			contentFilterRaw.Blocking = *contentFilter.Blocking
		}

		if contentFilter.Enabled != nil {
			contentFilterRaw.Enabled = *contentFilter.Enabled
		}

		if contentFilter.Name != nil {
			contentFilterRaw.Name = *contentFilter.Name
		}

		if contentFilter.Source != nil {
			contentFilterRaw.Source = string(*contentFilter.Source)
		}

		if contentFilter.SeverityThreshold != nil {
			contentFilterRaw.SeverityThreshold = string(*contentFilter.SeverityThreshold)
		}

		contentFilters = append(contentFilters, contentFilterRaw)
	}

	return contentFilters
}

func flattenCustomBlicklistsModel(inputList *[]raipolicies.CustomBlocklistConfig) []CustomBlocklistsModel {
	customBlocklists := make([]CustomBlocklistsModel, 0)

	if inputList == nil {
		return customBlocklists
	}

	for _, blockList := range *inputList {
		blockListRaw := CustomBlocklistsModel{}

		if blockList.Blocking != nil {
			blockListRaw.Blocking = *blockList.Blocking
		}

		if blockList.BlocklistName != nil {
			blockListRaw.BlocklistName = *blockList.BlocklistName
		}

		if blockList.Source != nil {
			blockListRaw.Source = string(*blockList.Source)
		}

		customBlocklists = append(customBlocklists, blockListRaw)
	}

	return customBlocklists
}
