// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cognitiveAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"kind": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"AIServices",
				"Academic",
				"AnomalyDetector",
				"Bing.Autosuggest",
				"Bing.Autosuggest.v7",
				"Bing.CustomSearch",
				"Bing.Search",
				"Bing.Search.v7",
				"Bing.Speech",
				"Bing.SpellCheck",
				"Bing.SpellCheck.v7",
				"CognitiveServices",
				"ComputerVision",
				"ContentModerator",
				"ConversationalLanguageUnderstanding",
				"ContentSafety",
				"CustomSpeech",
				"CustomVision.Prediction",
				"CustomVision.Training",
				"Emotion",
				"Face",
				"FormRecognizer",
				"ImmersiveReader",
				"LUIS",
				"LUIS.Authoring",
				"MetricsAdvisor",
				"OpenAI",
				"Personalizer",
				"QnAMaker",
				"Recommendations",
				"SpeakerRecognition",
				"Speech",
				"SpeechServices",
				"SpeechTranslation",
				"TextAnalytics",
				"TextTranslation",
				"WebLM",
			}, false),
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"C2", "C3", "C4", "D3", "DC0", "E0", "F0", "F1", "P0", "P1", "P2", "S", "S0", "S1", "S2", "S3", "S4", "S5", "S6",
			}, false),
		},

		"project_management_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"custom_subdomain_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},

					"identity_client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"dynamic_throttling_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"fqdns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"local_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"metrics_advisor_aad_client_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"metrics_advisor_aad_tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"metrics_advisor_super_user_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"metrics_advisor_website_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_acls": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			RequiredWith: []string{"custom_subdomain_name"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"default_action": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(cognitiveservicesaccounts.PossibleValuesForNetworkRuleAction(), false),
					},

					"ip_rules": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								validation.IsIPv4Address,
								validation.IsCIDRIPv4,
							),
						},
						Set: set.HashIPv4AddressOrCIDR,
					},

					"virtual_network_rules": {
						Type:       pluginsdk.TypeSet,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAuto,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"ignore_missing_vnet_service_endpoint": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},

					"bypass": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice(
							cognitiveservicesaccounts.PossibleValuesForByPassSelection(),
							false,
						),
					},
				},
			},
		},

		"network_injection": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"scenario": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(cognitiveservicesaccounts.ScenarioTypeAgent),
						}, false),
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"outbound_network_access_restricted": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"qna_runtime_endpoint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"custom_question_answering_search_service_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: services.ValidateSearchServiceID,
		},

		"custom_question_answering_search_service_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Sensitive:    true,
		},

		"storage": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_account_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"identity_client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}
