// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gallery_application

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.GalleryApplicationName,
		},

		"gallery_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSharedImageGalleryID,
		},

		"location": commonschema.Location(),

		"supported_os_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryapplications.PossibleValuesForOperatingSystemTypes(), false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"end_of_life_date": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"eula": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"privacy_statement_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"release_note_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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

			if oldVal, newVal := metadata.ResourceDiff.GetChange("privacy_statement_uri"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("privacy_statement_uri"); err != nil {
					return err
				}
			}

			if oldVal, newVal := metadata.ResourceDiff.GetChange("release_note_uri"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("release_note_uri"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
