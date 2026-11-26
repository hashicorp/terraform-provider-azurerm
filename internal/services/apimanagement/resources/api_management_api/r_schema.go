// Copyright IBM Corp. 2014, 2026
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
		"name": schemaz.SchemaApiManagementApiName("The name of the API Management API."),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"display_name": {
			Description:  "The display name of the API.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"path": {
			Description:  "The Path for this API Management API, which is a relative URL which uniquely identifies this API and all of its resource paths within the API Management Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.ApiManagementApiPath,
		},

		"protocols": {
			Description: "A list of protocols the operations in this API can be invoked.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(api.PossibleValuesForProtocol(), false),
			},
		},

		"revision": {
			Description:  "The Revision which used for this API.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"revision_description": {
			Description:  "The description of the API Revision of the API Management API.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Optional
		"api_type": {
			Description: "Type of API.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
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
						Description:  "The email address of the contact person/organization.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsEmailAddress,
					},
					"name": {
						Description:  "The name of the contact person/organization.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"url": {
						Description:  "Absolute URL of the contact information.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"description": {
			Description: "A description of the API Management API, which may include HTML formatting tags.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"import": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content_value": {
						Description:  "The Content from which the API Definition should be imported. When a `content_format` of `*-link-*` is specified this must be a URL, otherwise this must be defined inline. The URL must be accessible and return a valid document; otherwise, deployment may fail.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"content_format": {
						Description: "The format of the content from which the API Definition should be imported.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description: "This can only be specified when `content_format` is `wsdl` or `wsdl-link`.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_name": {
									Description:  "The name of service to import from WSDL.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"endpoint_name": {
									Description:  "The name of endpoint (port) to import from WSDL.",
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
						Description:  "The name of the license.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"url": {
						Description:  "Absolute URL of the license.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"service_url": {
			Description: "Absolute URL of the backend service implementing this API.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
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
						Description:  "The name of the HTTP Header which should be used for the Subscription Key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"query": {
						Description:  "The name of the QueryString parameter which should be used for the Subscription Key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"subscription_required": {
			Description: "Whether this API should require a subscription key.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"terms_of_service_url": {
			Description:  "Absolute URL of the Terms of Service for the API.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"source_api_id": {
			Description:  "The API ID of the source API, which could be in format `azurerm_api_management_api.example.id` or in format `azurerm_api_management_api.example.id;rev=1`.",
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
						Description:  "OAuth authorization server identifier. The name of an [OAuth2 Authorization Server](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/api_management_authorization_server).",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApiManagementChildName,
					},
					"scope": {
						Description: "Operations scope.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
						Description:  "OpenID Connect provider identifier. The name of an [OpenID Connect Provider](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/api_management_openid_connect_provider).",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApiManagementChildName,
					},
					"bearer_token_sending_methods": {
						Description: "How to send token to the server. A list of zero or more methods.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
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
			Description: "Whether this is the current API Revision.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"is_online": {
			Description: "Whether this API Revision is online/accessible via the Gateway.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"version": {
			Description: "The Version number of this API, if this API is versioned.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because the API assigns a version identifier when not explicitly set
			Computed: true,
		},

		"version_description": {
			Description:  "The description of the API Version of the API Management API.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_set_id": {
			Description: "The ID of the Version Set which this API is associated with.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because the API assigns a version set ID when not explicitly set
			Computed: true,
		},
	}
}
