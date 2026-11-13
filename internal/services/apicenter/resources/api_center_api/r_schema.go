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
			Type:         pluginsdk.TypeString,
			Description:  "The name of the API.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},

		"service_id": {
			Type:         pluginsdk.TypeString,
			Description:  "The ID of the API Center service.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: services.ValidateServiceID,
		},

		"kind": {
			Type:         pluginsdk.TypeString,
			Description:  "The kind of API.",
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
						Type:         pluginsdk.TypeString,
						Description:  "The email address of the API contact.",
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 100),
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Description:  "The name of the API contact.",
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 100),
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Description:  "The URL for the API contact.",
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},

		"custom_properties": helpers.CustomPropertiesSchema(),

		"description": {
			Type:         pluginsdk.TypeString,
			Description:  "A description of the API.",
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 1000),
		},

		"display_name": {
			Type:        pluginsdk.TypeString,
			Description: "The display name of the API.",
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
						Type:         pluginsdk.TypeString,
						Description:  "The URL of the external documentation.",
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"description": {
						Type:         pluginsdk.TypeString,
						Description:  "A description of the external documentation.",
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 500),
					},

					"title": {
						Type:         pluginsdk.TypeString,
						Description:  "The title of the external documentation.",
						Optional:     true,
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
						Type:          pluginsdk.TypeString,
						Description:   "The SPDX identifier for the API license.",
						Optional:      true,
						ValidateFunc:  validation.StringLenBetween(1, 50),
						ConflictsWith: []string{"license.0.url"},
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Description:  "The name of the API license.",
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 100),
					},

					"url": {
						Type:          pluginsdk.TypeString,
						Description:   "The URL of the API license.",
						Optional:      true,
						ValidateFunc:  validation.IsURLWithHTTPorHTTPS,
						ConflictsWith: []string{"license.0.identifier"},
					},
				},
			},
		},

		"summary": {
			Type:         pluginsdk.TypeString,
			Description:  "A short summary of the API.",
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 200),
		},

		"terms_of_service_url": {
			Type:         pluginsdk.TypeString,
			Description:  "The URL of the API's terms of service.",
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"lifecycle_stage": {
			Type:        pluginsdk.TypeString,
			Description: "The lifecycle stage of the API.",
			Computed:    true,
		},
	}
}
