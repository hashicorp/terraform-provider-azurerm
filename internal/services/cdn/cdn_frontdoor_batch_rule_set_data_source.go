// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	legacyrulesets "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = CdnFrontDoorBatchRuleSetDataSource{}

type CdnFrontDoorBatchRuleSetDataSource struct{}

type CdnFrontDoorBatchRuleSetDataSourceModel struct {
	Name                  string                           `tfschema:"name"`
	ProfileName           string                           `tfschema:"profile_name"`
	ResourceGroupName     string                           `tfschema:"resource_group_name"`
	BatchModeEnabled      bool                             `tfschema:"batch_mode_enabled"`
	CdnFrontDoorProfileID string                           `tfschema:"cdn_frontdoor_profile_id"`
	Rules                 []CdnFrontDoorBatchRuleRuleModel `tfschema:"rules"`
}

func (CdnFrontDoorBatchRuleSetDataSource) ResourceType() string {
	return "azurerm_cdn_frontdoor_batch_rule_set"
}

func (CdnFrontDoorBatchRuleSetDataSource) ModelObject() interface{} {
	return &CdnFrontDoorBatchRuleSetDataSourceModel{}
}

func (CdnFrontDoorBatchRuleSetDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.FrontDoorRuleSetName,
		},
		"profile_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.FrontDoorName,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (CdnFrontDoorBatchRuleSetDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"batch_mode_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"cdn_frontdoor_profile_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"rules": cdnFrontDoorBatchRulesComputedSchema(),
	}
}

func (CdnFrontDoorBatchRuleSetDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			batchModeRuleSetClient := metadata.Client.Cdn.FrontDoorRuleSetsClient_v2025_12_01

			var state CdnFrontDoorBatchRuleSetDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := rules.NewRuleSetID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.ProfileName, state.Name)
			ruleSetResourceId := legacyrulesets.NewRuleSetID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)

			resp, err := batchModeRuleSetClient.Get(ctx, ruleSetResourceId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.BatchMode == nil || !pointer.From(resp.Model.Properties.BatchMode) {
				return fmt.Errorf("retrieving %s: `batch_mode_enabled` must be `true` on the parent Rule Set", id)
			}

			flattenedState, err := flattenCdnFrontDoorBatchRuleSetDataSourceModel(&id, resp.Model)
			if err != nil {
				return err
			}

			state = flattenedState

			metadata.SetID(&id)
			return metadata.Encode(&state)
		},
	}
}

func cdnFrontDoorBatchRulesComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"behavior_on_match": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"order": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
			"actions":    cdnFrontDoorBatchRuleActionsComputedSchema(),
			"conditions": cdnFrontDoorBatchRuleConditionsComputedSchema(),
		}},
	}
}

func cdnFrontDoorBatchRuleActionsComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"url_redirect_action":                 cdnFrontDoorBatchRuleURLRedirectActionComputedSchema(),
			"url_rewrite_action":                  cdnFrontDoorBatchRuleURLRewriteActionComputedSchema(),
			"request_header_action":               cdnFrontDoorBatchRuleHeaderActionComputedSchema(),
			"response_header_action":              cdnFrontDoorBatchRuleHeaderActionComputedSchema(),
			"route_configuration_override_action": cdnFrontDoorBatchRuleRouteConfigurationOverrideActionComputedSchema(),
		}},
	}
}

func cdnFrontDoorBatchRuleURLRedirectActionComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"redirect_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"redirect_protocol": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"destination_path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"destination_hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"query_string": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"destination_fragment": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		}},
	}
}

func cdnFrontDoorBatchRuleURLRewriteActionComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"source_pattern": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"destination": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"preserve_unmatched_path": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		}},
	}
}

func cdnFrontDoorBatchRuleHeaderActionComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"header_action": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"header_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"value": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		}},
	}
}

func cdnFrontDoorBatchRuleRouteConfigurationOverrideActionComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"cdn_frontdoor_origin_group_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"forwarding_protocol": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"query_string_caching_behavior": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"query_string_parameters": batchComputedStringListSchema(),
			"compression_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
			"cache_behavior": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"cache_duration": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		}},
	}
}

func cdnFrontDoorBatchRuleConditionsComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"remote_address_condition":     batchConditionComputedSchema(batchComputedStringListSchema(), false, ""),
			"request_method_condition":     batchConditionComputedSchema(batchComputedStringSetSchema(), false, ""),
			"query_string_condition":       batchConditionComputedSchema(batchComputedStringListSchema(), true, ""),
			"post_args_condition":          batchConditionComputedSchema(batchComputedStringListSchema(), true, "post_args_name"),
			"request_uri_condition":        batchConditionComputedSchema(batchComputedStringListSchema(), true, ""),
			"request_header_condition":     batchConditionComputedSchema(batchComputedStringListSchema(), true, "header_name"),
			"request_body_condition":       batchConditionComputedSchema(batchComputedStringListSchema(), true, ""),
			"request_scheme_condition":     batchConditionComputedSchema(batchComputedStringListSchema(), false, ""),
			"url_path_condition":           batchConditionComputedSchema(batchComputedStringListSchema(), true, ""),
			"url_file_extension_condition": batchConditionComputedSchema(batchComputedStringListSchema(), true, ""),
			"url_filename_condition":       batchConditionComputedSchema(batchComputedStringListSchema(), true, ""),
			"http_version_condition":       batchConditionComputedSchema(batchComputedStringSetSchema(), false, ""),
			"cookies_condition":            batchConditionComputedSchema(batchComputedStringListSchema(), true, "cookie_name"),
			"is_device_condition":          batchConditionComputedSchema(batchComputedStringListSchema(), false, ""),
			"socket_address_condition":     batchConditionComputedSchema(batchComputedStringListSchema(), false, ""),
			"client_port_condition":        batchConditionComputedSchema(batchComputedStringListSchema(), false, ""),
			"server_port_condition":        batchConditionComputedSchema(batchComputedStringSetSchema(), false, ""),
			"host_name_condition":          batchConditionComputedSchema(batchComputedStringListSchema(), true, ""),
			"ssl_protocol_condition":       batchConditionComputedSchema(batchComputedStringSetSchema(), false, ""),
		}},
	}
}

func batchConditionComputedSchema(matchValuesSchema *pluginsdk.Schema, includeTransforms bool, selectorField string) *pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"operator": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"negate_condition": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"match_values": matchValuesSchema,
	}

	if includeTransforms {
		schema["transforms"] = batchComputedStringSetSchema()
	}

	if selectorField != "" {
		schema[selectorField] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Computed: true,
		}
	}

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem:     &pluginsdk.Resource{Schema: schema},
	}
}

func batchComputedStringListSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func batchComputedStringSetSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Computed: true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func flattenCdnFrontDoorBatchRuleSetDataSourceModel(id *rules.RuleSetId, model *azuresdkhacks.BatchRuleSetResource) (CdnFrontDoorBatchRuleSetDataSourceModel, error) {
	state := CdnFrontDoorBatchRuleSetDataSourceModel{
		Name:                  id.RuleSetName,
		ProfileName:           id.ProfileName,
		ResourceGroupName:     id.ResourceGroupName,
		BatchModeEnabled:      true,
		CdnFrontDoorProfileID: profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID(),
	}

	if model != nil && model.Properties != nil && model.Properties.Rules != nil {
		rulesState, err := flattenCdnFrontDoorBatchRulesDataSource(pointer.From(model.Properties.Rules))
		if err != nil {
			return CdnFrontDoorBatchRuleSetDataSourceModel{}, fmt.Errorf("flattening `rules`: %+v", err)
		}
		state.Rules = rulesState
	}

	return state, nil
}

func flattenCdnFrontDoorBatchRulesDataSource(input []azuresdkhacks.BatchRuleProperties) ([]CdnFrontDoorBatchRuleRuleModel, error) {
	results := make([]CdnFrontDoorBatchRuleRuleModel, 0, len(input))
	sorted := append([]azuresdkhacks.BatchRuleProperties(nil), input...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return batchRuleOrder(sorted[i]) < batchRuleOrder(sorted[j])
	})

	for _, item := range sorted {
		ruleState, err := flattenCdnFrontDoorBatchRuleDataSource(item)
		if err != nil {
			return []CdnFrontDoorBatchRuleRuleModel{}, err
		}
		results = append(results, ruleState)
	}

	return results, nil
}

func flattenCdnFrontDoorBatchRuleDataSource(input azuresdkhacks.BatchRuleProperties) (CdnFrontDoorBatchRuleRuleModel, error) {
	state := CdnFrontDoorBatchRuleRuleModel{
		Name: pointer.From(input.Name),
	}

	state.BehaviorOnMatch = string(pointer.From(input.MatchProcessingBehavior))
	state.Order = pointer.From(input.Order)

	actions, err := flattenCdnFrontDoorBatchRuleActions(input.Actions)
	if err != nil {
		return CdnFrontDoorBatchRuleRuleModel{}, fmt.Errorf("flattening `actions`: %+v", err)
	}
	state.Actions = actions

	conditions, err := flattenCdnFrontDoorBatchRuleConditions(input.Conditions)
	if err != nil {
		return CdnFrontDoorBatchRuleRuleModel{}, fmt.Errorf("flattening `conditions`: %+v", err)
	}
	state.Conditions = conditions

	return state, nil
}
