// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_app_attach_package

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/appattachpackage"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// `ValidateFunc` coded according to portal
			ValidateFunc: validation.StringDoesNotContainAny("\\/+?&"),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"host_pool_ids": {
			// `TypeSet` because order is not guaranteed
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: hostpool.ValidateHostPoolID,
			},
		},

		"msix_package_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_share_file_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageShareDataPlaneID,
		},

		"health_check_status_on_failure": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      appattachpackage.FailHealthCheckOnStagingFailureNeedsAssistance,
			ValidateFunc: validation.StringInSlice(appattachpackage.PossibleValuesForFailHealthCheckOnStagingFailure(), false),
		},

		"register_at_log_on_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"state_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"last_updated": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"package_applications": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"app_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"app_user_model_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"friendly_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"icon_image_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"raw_icon": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"raw_png": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"package_family_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"package_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"package_relative_path": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			oldHostPoolIds, newHostPoolIds := metadata.ResourceDiff.GetChange("host_pool_ids")

			// If number of `host_pool_ids` elements has changed, actual elements will be updated incorrectly to be different from Terraform configurations
			if oldHostPoolIds.(*pluginsdk.Set).Len() != newHostPoolIds.(*pluginsdk.Set).Len() {
				metadata.ResourceDiff.ForceNew("host_pool_ids")
			}

			return nil
		},
	}
}
