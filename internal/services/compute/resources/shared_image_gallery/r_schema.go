// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package shared_image_gallery

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sharedImageGallerySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Shared Image Gallery.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SharedImageGalleryName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"description": {
			Description: "A description for this Shared Image Gallery.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"sharing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"permission": {
						Description:  "The permission of the Shared Image Gallery when sharing.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(galleries.PossibleValuesForGallerySharingPermissionTypes(), false),
					},

					"community_gallery": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"eula": {
									Description:  "The End User Licence Agreement for the Shared Image Gallery.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"prefix": {
									Description:  "Prefix of the community public name for the Shared Image Gallery.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validate.SharedImageGalleryPrefix,
								},
								"publisher_email": {
									Description:  "Email of the publisher for the Shared Image Gallery.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"publisher_uri": {
									Description:  "URI of the publisher for the Shared Image Gallery.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"name": {
									Description: "The community public name of the Shared Image Gallery.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"unique_name": {
			Description: "The Unique Name for this Shared Image Gallery.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
