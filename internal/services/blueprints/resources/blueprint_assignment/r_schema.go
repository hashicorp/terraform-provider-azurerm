// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Blueprint Assignment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_subscription_id": {
			Description:  "The Subscription ID the Blueprint Published Version is to be applied to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"location": commonschema.Location(),

		"identity": commonschema.SystemOrUserAssignedIdentityRequired(),

		"version_id": {
			Description:  "The ID of the Published Version of the blueprint to be assigned.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: publishedblueprint.ValidateScopedVersionID,
		},

		"parameter_values": {
			Description:      "A JSON string to supply Blueprint Assignment parameter values.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			StateFunc:        normalizeAssignmentParameterValuesJSON,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"resource_groups": {
			Description:      "A JSON string to supply the Blueprint Resource Group information.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			StateFunc:        normalizeAssignmentResourceGroupValuesJSON,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"lock_mode": {
			Description:  "The locking mode of the Blueprint Assignment. One of `None` (Default), `AllResourcesReadOnly`, or `AllResourcesDoNotDelete`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(assignment.AssignmentLockModeNone),
			ValidateFunc: validation.StringInSlice(assignment.PossibleValuesForAssignmentLockMode(), false),
			// The first character of value returned by the service is always in lower case.
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"lock_exclude_principals": {
			Description: "A list of up to 5 Principal IDs that are permitted to bypass the locks applied by the Blueprint.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    5,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"lock_exclude_actions": {
			Description: "A list of up to 200 actions that are permitted to bypass the locks applied by the Blueprint.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    200,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Description: "The Description on the Blueprint.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The display name of the blueprint.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"blueprint_name": {
			Description: "The name of the blueprint assigned.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
