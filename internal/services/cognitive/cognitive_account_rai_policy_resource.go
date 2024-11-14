package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = &CognitiveAccountRaiPolicyResource{}

type CognitiveAccountRaiPolicyResource struct{}

type AccountRaiPolicyContentFilter struct {
	Name              string `tfschema:"name"`
	FilterEnabled     bool   `tfschema:"filter_enabled"`
	BlockEnabled      bool   `tfschema:"block_enabled"`
	SeverityThreshold string `tfschema:"severity_threshold"`
	Source            string `tfschema:"source"`
}

type AccountRaiPolicyCustomBlock struct {
	Name         string `tfschema:"name"`
	BlockEnabled bool   `tfschema:"block_enabled"`
	Source       string `tfschema:"source"`
}

type AccountRaiPolicyResourceModel struct {
	Name            string                          `tfschema:"name"`
	AccountId       string                          `tfschema:"cognitive_account_id"`
	BasePolicyName  string                          `tfschema:"base_policy_name"`
	ContentFilter   []AccountRaiPolicyContentFilter `tfschema:"content_filter"`
	CustomBlockList []AccountRaiPolicyCustomBlock   `tfschema:"custom_blocklist"`
	Mode            string                          `tfschema:"mode"`
	Tags            map[string]string               `tfschema:"tags"`
	Type            string                          `tfschema:"type"`
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

		"custom_blocklist": {
			Type:     pluginsdk.TypeList,
			Optional: true,
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
	return map[string]*pluginsdk.Schema{
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
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
				return err
			}

			cognitiveAccountId, err := raipolicies.ParseAccountID(model.AccountId)
			if err != nil {
				return err
			}

			id := raipolicies.NewRaiPolicyID(subscriptionId, cognitiveAccountId.ResourceGroupName, cognitiveAccountId.AccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			locks.ByID(cognitiveAccountId.ID())
			defer locks.UnlockByID(cognitiveAccountId.ID())

			raiPolicy := raipolicies.RaiPolicy{
				Name: pointer.To(model.Name),
				Properties: &raipolicies.RaiPolicyProperties{
					BasePolicyName:   pointer.To(model.BasePolicyName),
					ContentFilters:   expandRaiPolicyContentFilters(model.ContentFilter),
					CustomBlocklists: expandCustomBlockLists(model.CustomBlockList),
					Mode:             pointer.To(raipolicies.RaiPolicyMode(model.Mode)),
				},
				Tags: pointer.To(model.Tags),
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

			state := AccountRaiPolicyResourceModel{}

			if model := resp.Model; model != nil {
				state.Name = pointer.From(model.Name)
				state.AccountId = cognitiveAccountId.ID()
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.BasePolicyName = pointer.From(props.BasePolicyName)
					state.ContentFilter = flattenRaiPolicyContentFilters(props.ContentFilters)
					state.CustomBlockList = flattenRaiCustomBlockLists(props.CustomBlocklists)
					state.Mode = string(pointer.From(props.Mode))
					state.Type = string(pointer.From(props.Type))
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

			id, err := raipolicies.ParseRaiPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AccountRaiPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			cognitiveAccountId := raipolicies.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(cognitiveAccountId.ID())
			defer locks.UnlockByID(cognitiveAccountId.ID())

			props := resp.Model

			if metadata.ResourceData.HasChange("content_filter") {
				props.Properties.ContentFilters = expandRaiPolicyContentFilters(model.ContentFilter)
			}

			if metadata.ResourceData.HasChange("custom_blocklist") {
				props.Properties.CustomBlocklists = expandCustomBlockLists(model.CustomBlockList)
			}

			if metadata.ResourceData.HasChange("mode") {
				props.Properties.Mode = pointer.To(raipolicies.RaiPolicyMode(model.Mode))
			}

			if metadata.ResourceData.HasChange("tags") {
				props.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *props); err != nil {
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

func expandCustomBlockLists(list []AccountRaiPolicyCustomBlock) *[]raipolicies.CustomBlocklistConfig {
	if list == nil {
		return nil
	}

	customBlockLists := make([]raipolicies.CustomBlocklistConfig, 0, len(list))
	for _, block := range list {
		customBlockLists = append(customBlockLists, raipolicies.CustomBlocklistConfig{
			BlocklistName: pointer.To(block.Name),
			Blocking:      pointer.To(block.BlockEnabled),
			Source:        pointer.To(raipolicies.RaiPolicyContentSource(block.Source)),
		})
	}
	return &customBlockLists
}

func flattenRaiCustomBlockLists(blocklists *[]raipolicies.CustomBlocklistConfig) []AccountRaiPolicyCustomBlock {
	if blocklists == nil {
		return nil
	}

	customBlockLists := make([]AccountRaiPolicyCustomBlock, 0, len(*blocklists))
	for _, block := range *blocklists {
		customBlockLists = append(customBlockLists, AccountRaiPolicyCustomBlock{
			Name:         pointer.From(block.BlocklistName),
			BlockEnabled: pointer.From(block.Blocking),
			Source:       string(pointer.From(block.Source)),
		})
	}
	return customBlockLists
}

func flattenRaiPolicyContentFilters(filters *[]raipolicies.RaiPolicyContentFilter) []AccountRaiPolicyContentFilter {
	if filters == nil {
		return nil
	}

	contentFilters := make([]AccountRaiPolicyContentFilter, 0, len(*filters))
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
