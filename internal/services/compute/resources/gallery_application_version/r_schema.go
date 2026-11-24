// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package gallery_application_version

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The version name of the Gallery Application Version, such as `1.0.0`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.GalleryApplicationVersionName,
		},

		"gallery_application_id": {
			Description:  "The ID of the Gallery Application.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: galleryapplications.ValidateApplicationID,
		},

		"location": commonschema.Location(),

		"config_file": {
			Description:  "The name of the config file on the VM.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enable_health_check": {
			Description: "Whether the Gallery Application should report health.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"end_of_life_date": {
			Description:      "The end of life date in RFC3339 format of the Gallery Application Version.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"exclude_from_latest": {
			Description: "Whether the Gallery Application Version should be excluded from the `latest` filter. If set to `true` this Gallery Application Version won't be returned for the `latest` version.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"manage_action": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"install": {
						Description:  "The command to install the Gallery Application.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 4096),
					},

					"remove": {
						Description:  "The command to remove the Gallery Application.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 4096),
					},

					"update": {
						Description:  "The command to update the Gallery Application.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringLenBetween(1, 4096),
					},
				},
			},
		},

		"package_file": {
			Description:  "The name of the package file on the VM.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"source": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"media_link": {
						Description:  "The Storage Blob URI of the source application package.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},

					"default_configuration_link": {
						Description:  "The Storage Blob URI of the default configuration.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},
				},
			},
		},

		"target_region": {
			// This needs to be a `TypeList` due to the `StateFunc` on the nested property `name`
			// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/160
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": commonschema.LocationWithoutForceNew(),

					"regional_replica_count": {
						Description:  "The number of replicas of the Gallery Application Version to be created per region.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 10),
					},

					"exclude_from_latest": {
						Description: "Whether this Gallery Application Version should be excluded from the `latest` filter. If set to `true`, this Gallery Application Version won't be returned for the `latest` version.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"storage_account_type": {
						Description:  "The Storage Account type for the Gallery Application Version.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(galleryapplicationversions.PossibleValuesForStorageAccountType(), false),
						Default:      string(galleryapplicationversions.StorageAccountTypeStandardLRS),
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("end_of_life_date"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("end_of_life_date"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
