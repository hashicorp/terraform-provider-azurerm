// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_definition

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the API definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},

		"api_center_api_version_id": {
			Description:  "The ID of the API version to which this API definition belongs.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apiversions.ValidateVersionID,
		},

		"description": {
			Description:  "The description of the API definition.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Description: "The display name of the API definition.",
			Type:        pluginsdk.TypeString,
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
						Description: "The name of the Specification.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"openapi", "wsdl", "wadl", "graphql", "grpc", "asyncapi", "raml",
							"other",
						}, false),
					},

					"content": {
						Description: "The content of the Specification.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						// Note: O+C if URL is set, Azure will be used to fetch the specification content and only return that
						Computed:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"specification.0.url"},
						AtLeastOneOf:  []string{"specification.0.content", "specification.0.url"},
					},

					"url": {
						Description:   "The URL containing the Specification content.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.IsURLWithHTTPS,
						ConflictsWith: []string{"specification.0.content"},
						AtLeastOneOf:  []string{"specification.0.content", "specification.0.url"},
					},

					"version": {
						Description: "The version of the Specification.",
						Type:        pluginsdk.TypeString,
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

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(_ context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff
			if diff.HasChange("specification.#") && diff.NewValueKnown("specification.#") && diff.Get("specification.#").(int) == 0 {
				return diff.ForceNew("specification")
			}
			return nil
		},
	}
}
