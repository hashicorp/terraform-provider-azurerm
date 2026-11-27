// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the API.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},

		"api_center_service_id": {
			Description:  "The ID of the API Center service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: services.ValidateServiceID,
		},

		"api_type": {
			Description:  "The type of API.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(apis.PossibleValuesForApiKind(), false),
		},

		"contact": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Description:  "The email address of the API contact.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 100),
					},

					"name": {
						Description:  "The name of the API contact.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 100),
					},

					"url": {
						Description:  "The URL for the API contact.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"custom_properties_json": helpers.CustomPropertiesSchema(),

		"description": {
			Description:  "A description of the API.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 1000),
		},

		"display_name": {
			Description: "The display name of the API.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C if display name is not set, it will default to the name of the API
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"external_documentation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 20,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url": {
						Description:  "The URL of the external documentation.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"description": {
						Description:  "A description of the external documentation.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 500),
					},

					"title": {
						Description:  "The title of the external documentation.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 50),
					},
				},
			},
		},

		"license": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identifier": {
						Description:   "The SPDX identifier for the API license.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringLenBetween(1, 50),
						ConflictsWith: []string{"license.0.url"},
					},

					"name": {
						Description:  "The name of the API license.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 100),
					},

					"url": {
						Description:   "The URL of the API license.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.IsURLWithHTTPorHTTPS,
						ConflictsWith: []string{"license.0.identifier"},
					},
				},
			},
		},

		"summary": {
			Description:  "A short summary of the API.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 200),
		},

		"terms_of_service_url": {
			Description:  "The URL of the API's terms of service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"lifecycle_stage": {
			Description: "The lifecycle stage of the API.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
