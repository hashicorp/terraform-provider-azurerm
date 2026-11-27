// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_metadata_schema

import (
	"context"
	"errors"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/metadataschemas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apicenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the metadata schema.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Name(),
		},

		"api_center_service_id": {
			Description:  "The ID of the API Center service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: services.ValidateServiceID,
		},

		"metadata_assignment": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"entity": {
						Description:  "The type of entity to which the metadata applies.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(metadataschemas.PossibleValuesForMetadataAssignmentEntity(), false),
					},

					"deprecated": {
						Description: "Whether the metadata assignment is deprecated.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"required": {
						Description: "Whether the metadata assignment is required.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
				},
			},
		},

		"schema_json": {
			Description:  "The JSON schema defining the metadata.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsJSON,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			if !metadata.ResourceDiff.HasChange("metadata_assignment") {
				return nil
			}

			var model ApiCenterMetadataSchemaModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return err
			}

			for _, assignment := range model.MetadataAssignment {
				if assignment.Deprecated && assignment.Required {
					return errors.New("metadata_assignment cannot be both deprecated and required")
				}
			}

			return nil
		},
	}
}
