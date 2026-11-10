// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesprojects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}
			// Azure API issue: description and display_name cannot be updated to empty values, see: https://github.com/Azure/azure-rest-api-specs/issues/38422
			if oldDesc, newDesc := metadata.ResourceDiff.GetChange("description"); oldDesc.(string) != "" && newDesc.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("description"); err != nil {
					return err
				}
			}

			if oldDisplay, newDisplay := metadata.ResourceDiff.GetChange("display_name"); oldDisplay.(string) != "" && newDisplay.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("display_name"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountProjectName(),
		},

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&cognitiveservicesprojects.AccountId{}),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoints": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"default": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}
