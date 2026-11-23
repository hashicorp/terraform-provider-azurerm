// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Cognitive Service Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"kind": {
			Description: "The type of Cognitive Service Account that should be created. More information on [upgrade and rollback scenario](https://learn.microsoft.com/azure/ai-foundry/how-to/upgrade-azure-openai?tabs=portal).",
			Type:        pluginsdk.TypeString,
			Required:    true,
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
			Description: "The SKU Name for this Cognitive Service Account.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"C2", "C3", "C4", "D3", "DC0", "E0", "F0", "F1", "P0", "P1", "P2", "S", "S0", "S1", "S2", "S3", "S4", "S5", "S6",
			}, false),
		},

		"project_management_enabled": {
			Description: "Whether project management is enabled. Can only be set to `true` when `kind` is set to `AIServices`. Once enabled, disabling `project_management_enabled` forces a new resource to be created unless `kind` is set to `OpenAI`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"custom_subdomain_name": {
			Description:  "The subdomain name used for Entra ID token-based authentication. This attribute is required when `network_acls` is specified. This attribute is also required when using the OpenAI service with libraries which assume the Azure OpenAI endpoint is a subdomain on `https://openai.azure.com/`, e.g. `https://<custom_subdomain_name>.openai.azure.com/`. This can be specified during creation or added later, but once set changing this forces a new resource to be created.",
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
						Description:  "The ID of the Key Vault Key which should be used to Encrypt the data in this Cognitive Account.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},

					"identity_client_id": {
						Description:  "The Client ID of the User Assigned Identity that has access to the key. This property only needs to be specified when there're multiple identities attached to the Cognitive Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"dynamic_throttling_enabled": {
			Description: "Whether to enable the dynamic throttling for this Cognitive Service Account. This attribute cannot be set when the `kind` is `OpenAI` or `AIServices`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"fqdns": {
			Description: "List of FQDNs allowed for the Cognitive Account.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"local_auth_enabled": {
			Description: "Whether local authentication methods is enabled for the Cognitive Account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"metrics_advisor_aad_client_id": {
			Description:  "The Azure AD Client ID (Application ID). This attribute is only set when kind is `MetricsAdvisor`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"metrics_advisor_aad_tenant_id": {
			Description:  "The Azure AD Tenant ID. This attribute is only set when kind is `MetricsAdvisor`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"metrics_advisor_super_user_name": {
			Description:  "The super user of Metrics Advisor. This attribute is only set when kind is `MetricsAdvisor`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"metrics_advisor_website_name": {
			Description:  "The website name of Metrics Advisor. This attribute is only set when kind is `MetricsAdvisor`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_acls": {
			Description:  "When this property is specified, `custom_subdomain_name` is also required to be set.",
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
						Description: "One or more IP Addresses, or CIDR Blocks which should be able to access the Cognitive Account.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
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
									Description: "The ID of the subnet which should be able to access this Cognitive Account.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},

								"ignore_missing_vnet_service_endpoint": {
									Description: "Whether ignore missing VNet service endpoint or not.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},
							},
						},
					},

					"bypass": {
						Description: "Whether to allow trusted Azure Services to access the service.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice(
							cognitiveservicesaccounts.PossibleValuesForByPassSelection(),
							false,
						),
					},
				},
			},
		},

		"network_injection": {
			Description: "Only applicable if the `kind` is set to `AIServices`.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"scenario": {
						Description: "What features network injection applies to.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(cognitiveservicesaccounts.ScenarioTypeAgent),
						}, false),
					},

					"subnet_id": {
						Description:  "The ID of the subnet which the Agent Client is injected into.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"outbound_network_access_restricted": {
			Description: "Whether outbound network access is restricted for the Cognitive Account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is allowed for the Cognitive Account.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"qna_runtime_endpoint": {
			Description:  "A URL to link a QnAMaker Cognitive Account to a QnA runtime.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"custom_question_answering_search_service_id": {
			Description:  "If `kind` is `TextAnalytics` this specifies the ID of the Search Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: services.ValidateSearchServiceID,
		},

		"custom_question_answering_search_service_key": {
			Description:  "If `kind` is `TextAnalytics` this specifies the key of the Search Service.",
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
						Description:  "Full resource ID of a Microsoft.Storage resource.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"identity_client_id": {
						Description:  "The client ID of the Managed Identity associated with the storage resource.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"endpoint": {
			Description: "The endpoint used to connect to the Cognitive Service Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"primary_access_key": {
			Description: "A primary access key which can be used to connect to the Cognitive Service Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_access_key": {
			Description: "The secondary access key which can be used to connect to the Cognitive Service Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
	}
}
