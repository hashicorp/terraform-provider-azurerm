package migration

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ pluginsdk.StateUpgrade = ApiV0ToV1{}

type ApiV0ToV1 struct{}

func (ApiV0ToV1) Schema() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementApiName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.ApiManagementApiPath,
		},

		"protocols": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.ProtocolHTTP),
					string(apimanagement.ProtocolHTTPS),
					string(apimanagement.ProtocolWs),
					string(apimanagement.ProtocolWss),
				}, false),
			},
		},

		"revision": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"revision_description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Optional
		"api_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(apimanagement.APITypeGraphql),
				string(apimanagement.APITypeHTTP),
				string(apimanagement.APITypeSoap),
				string(apimanagement.APITypeWebsocket),
			}, false),
		},

		"contact": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.EmailAddress,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"import": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content_value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"content_format": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(apimanagement.ContentFormatOpenapi),
							string(apimanagement.ContentFormatOpenapijson),
							string(apimanagement.ContentFormatOpenapijsonLink),
							string(apimanagement.ContentFormatOpenapiLink),
							string(apimanagement.ContentFormatSwaggerJSON),
							string(apimanagement.ContentFormatSwaggerLinkJSON),
							string(apimanagement.ContentFormatWadlLinkJSON),
							string(apimanagement.ContentFormatWadlXML),
							string(apimanagement.ContentFormatWsdl),
							string(apimanagement.ContentFormatWsdlLink),
						}, false),
					},

					"wsdl_selector": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"endpoint_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"license": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"service_url": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"subscription_key_parameter_names": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"header": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"query": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"subscription_required": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"terms_of_service_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"source_api_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ApiID,
		},

		"oauth2_authorization": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorization_server_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApiManagementChildName,
					},
					"scope": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// There is currently no validation, as any length and characters can be used in the field
					},
				},
			},
		},

		"openid_authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"openid_provider_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApiManagementChildName,
					},
					"bearer_token_sending_methods": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(apimanagement.BearerTokenSendingMethodsAuthorizationHeader),
								string(apimanagement.BearerTokenSendingMethodsQuery),
							}, false),
						},
					},
				},
			},
		},

		// Computed
		"is_current": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"is_online": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			Optional: true,
		},

		"version_description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_set_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			Optional: true,
		},
	}

	if !features.FourPointOhBeta() {
		schema["api_type"].ConflictsWith = []string{"soap_pass_through"}

		schema["soap_pass_through"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			Deprecated:    "`soap_pass_through` will be removed in favour of the property `api_type` in version 4.0 of the AzureRM Provider",
			ConflictsWith: []string{"api_type"},
		}
	}

	return schema
}

func (ApiV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		rawState["id"] = rawState["id"].(string) + ";rev=" + rawState["revision"].(string)
		return rawState, nil
	}
}
