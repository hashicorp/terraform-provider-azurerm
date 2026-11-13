// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_definition

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Description:  "The name of the API definition.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},

		"api_version_id": {
			Type:         pluginsdk.TypeString,
			Description:  "The ID of the API version to which this API definition belongs.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apiversions.ValidateVersionID,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Description:  "The description of the API definition.",
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:        pluginsdk.TypeString,
			Description: "The display name of the API definition.",
			Optional:    true,
			// Note: O+C if display name is not set, it will default to the name of the API definition
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"specification": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:        pluginsdk.TypeString,
						Description: "The name of the Specification.",
						Required:    true,
					},

					"content": {
						Type:        pluginsdk.TypeString,
						Description: "The content of the Specification.",
						Optional:    true,
						AtLeastOneOf: []string{
							"specification.0.url",
							"specification.0.content",
						},
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"url": {
						Type:        pluginsdk.TypeString,
						Description: "The URL containing the Specification content.",
						Optional:    true,
						// Note: O+C if URL is set, Azure will be used to fetch the specification content and only return that
						Computed:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
						AtLeastOneOf: []string{
							"specification.0.url",
							"specification.0.content",
						},
					},

					"version": {
						Type:        pluginsdk.TypeString,
						Description: "The version of the Specification.",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
