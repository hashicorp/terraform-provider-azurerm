// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "Specifies the name of this Workbook as a UUID/GUID. It should not contain any uppercase letters.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.IsUUID,
				validate.StringDoesNotContainUpperCaseLetter,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Description:  "Specifies the user-defined name (display name) of the workbook.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_json": {
			Description:      "Configuration of this particular workbook. Configuration data is a string containing valid JSON.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"source_id": {
			Description:  "Resource ID for a source resource. It should not contain any uppercase letters.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "azure monitor",
			ValidateFunc: validate.StringDoesNotContainUpperCaseLetter,
		},

		"category": {
			Description:  "Workbook category, as defined by the user at creation time. There may be additional category types beyond the following: `workbook`, `sentinel`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "workbook",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Description:  "Specifies the description of the workbook.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

		"storage_container_id": {
			Description:  "Specifies the Resource Manager ID of the Storage Container when bring your own storage is used.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{
				"identity",
			},
		},

		"tags": {
			Description:  "A mapping of tags which should be assigned to the Workbook.",
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ValidateFunc: validate.WorkbookTags,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
