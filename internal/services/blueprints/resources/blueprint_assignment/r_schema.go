// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package blueprint_assignment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/publishedblueprint"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func blueprintAssignmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"location": commonschema.Location(),

		"identity": commonschema.SystemOrUserAssignedIdentityRequired(),

		"version_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: publishedblueprint.ValidateScopedVersionID,
		},

		"parameter_values": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			StateFunc:        normalizeAssignmentParameterValuesJSON,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"resource_groups": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			StateFunc:        normalizeAssignmentResourceGroupValuesJSON,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"lock_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(assignment.AssignmentLockModeNone),
			ValidateFunc: validation.StringInSlice(assignment.PossibleValuesForAssignmentLockMode(), false),
			// The first character of value returned by the service is always in lower case.
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"lock_exclude_principals": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 5,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"lock_exclude_actions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 200,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"blueprint_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
