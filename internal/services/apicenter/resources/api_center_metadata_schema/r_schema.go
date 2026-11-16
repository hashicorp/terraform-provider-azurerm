// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_metadata_schema

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/metadataschemas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Description:  "The name of the metadata schema.",
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"service_id": {
			Type:         pluginsdk.TypeString,
			Description:  "The ID of the API Center service.",
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
						Type:         pluginsdk.TypeString,
						Description:  "The type of entity to which the metadata applies.",
						Required:     true,
						ValidateFunc: validation.StringInSlice(metadataschemas.PossibleValuesForMetadataAssignmentEntity(), false),
					},

					"deprecated": {
						Type:        pluginsdk.TypeBool,
						Description: "Whether the metadata assignment is deprecated.",
						Optional:    true,
					},

					"required": {
						Type:        pluginsdk.TypeBool,
						Description: "Whether the metadata assignment is required.",
						Optional:    true,
					},
				},
			},
		},

		"schema": {
			Type:         pluginsdk.TypeString,
			Description:  "The JSON schema defining the metadata.",
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
			rd := metadata.ResourceDiff
			if rd.Id() == "" || !rd.HasChange("metadata_assignment") {
				return nil
			}

			oldRaw, newRaw := rd.GetChange("metadata_assignment")
			oldAssignments := oldRaw.([]interface{})
			newAssignments := newRaw.([]interface{})

			for _, oldRawAssignment := range oldAssignments {
				oldAssignment := oldRawAssignment.(map[string]interface{})
				for _, newRawAssignment := range newAssignments {
					newAssignment := newRawAssignment.(map[string]interface{})
					if oldAssignment["entity"] == newAssignment["entity"] &&
						oldAssignment["required"] != true && newAssignment["required"] == true {
						return rd.ForceNew("metadata_assignment")
					}
				}
			}

			return nil
		},
	}
}
