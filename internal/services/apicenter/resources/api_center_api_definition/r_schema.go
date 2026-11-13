package api_center_api_definition

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_version_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apiversions.ValidateVersionID,
		},
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},
		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// O+C if display name is not set, it will default to the name of the API definition
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
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						// O+C if URL is set, Azure will be used to fetch the specification content and only return that
						Computed:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
						AtLeastOneOf: []string{
							"specification.0.url",
							"specification.0.content",
						},
					},
					"content": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"specification.0.url",
							"specification.0.content",
						},
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
