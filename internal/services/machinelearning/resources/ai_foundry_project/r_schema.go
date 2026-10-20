// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ai_foundry_project

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_-]{2,32}$"),
				"AI Services Project name must be 2 - 32 characters long, contain only letters, numbers and hyphens.",
			),
		},

		"location": commonschema.Location(),

		"ai_services_hub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"high_business_impact_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			// NOTE: O+C creating a project that has encryption enabled with system assigned identity will set this property to true
			Computed: true,
			ForceNew: true,
		},

		"primary_user_assigned_identity": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			RequiredWith: []string{"identity"},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"project_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
