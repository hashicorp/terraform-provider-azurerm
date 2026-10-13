// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementApiSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementApiName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.ApiManagementApiPath,
		},

		"protocols": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(api.PossibleValuesForProtocol(), false),
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
			// Note: O+C because the API sets a default api_type (http) when not specified
			Computed:     true,
			ValidateFunc: validation.StringInSlice(api.PossibleValuesForApiType(), false),
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
						ValidateFunc: validation.IsEmailAddress,
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
							string(api.ContentFormatOpenapi),
							string(api.ContentFormatOpenapiPositivejson),
							string(api.ContentFormatOpenapiPositivejsonNegativelink),
							string(api.ContentFormatOpenapiNegativelink),
							string(api.ContentFormatSwaggerNegativejson),
							string(api.ContentFormatSwaggerNegativelinkNegativejson),
							string(api.ContentFormatWadlNegativelinkNegativejson),
							string(api.ContentFormatWadlNegativexml),
							string(api.ContentFormatWsdl),
							string(api.ContentFormatWsdlNegativelink),
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
			// Note: O+C because Azure returns a computed service_url when not specified
			Computed: true,
		},

		"subscription_key_parameter_names": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// Note: O+C because the API returns default subscription key parameter names when not specified
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
			ValidateFunc: validation.AsGeneratedID(api.ParseApiIDInsensitively),
		},

		"oauth2_authorization": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"openid_authentication"},
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
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"oauth2_authorization"},
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
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(api.PossibleValuesForBearerTokenSendingMethods(), false),
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
			// Note: O+C because the API assigns a version identifier when not explicitly set
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
			// Note: O+C because the API assigns a version set ID when not explicitly set
			Optional: true,
		},
	}
}
