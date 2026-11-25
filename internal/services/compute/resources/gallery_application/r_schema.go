// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Gallery Application.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.GalleryApplicationName,
		},

		"gallery_id": {
			Description:  "The ID of the Shared Image Gallery.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSharedImageGalleryID,
		},

		"location": commonschema.Location(),

		"supported_os_type": {
			Description:  "The type of the Operating System supported for the Gallery Application.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryapplications.PossibleValuesForOperatingSystemTypes(), false),
		},

		"description": {
			Description: "A description of the Gallery Application.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"end_of_life_date": {
			Description:      "The end of life date in RFC3339 format of the Gallery Application.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"eula": {
			Description:  "The End User Licence Agreement of the Gallery Application.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"privacy_statement_uri": {
			Description:  "The URI containing the Privacy Statement associated with the Gallery Application.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"release_note_uri": {
			Description:  "The URI containing the Release Notes associated with the Gallery Application.",
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
