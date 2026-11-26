// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project_environment_type

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Dev Center Project Environment Type.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterProjectEnvironmentTypeName,
		},

		"location": commonschema.Location(),

		"dev_center_project_id": commonschema.ResourceIDReferenceRequiredForceNew(&environmenttypes.ProjectId{}, "Dev Center Project"),

		"deployment_target_id": {
			Description:  "The ID of the subscription that the Environment Type will be mapped to. The environment's resources will be deployed into this subscription.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"creator_role_assignment_roles": {
			Description: "A list of roles to assign to the environment creator.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"user_role_assignment": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"user_id": {
						Description:  "The user object ID that is assigned roles.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},

					"roles": {
						Description: "A list of roles to assign to the `user_id`.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
