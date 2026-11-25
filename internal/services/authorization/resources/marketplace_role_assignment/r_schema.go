// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package marketplace_role_assignment

import (
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_id": {
			Description: "The ID of the Principal (User, Group or Service Principal) to assign the Role Definition to.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"name": {
			Description:  "A unique UUID/GUID for this Role Assignment - one will be generated if not specified.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"role_definition_id": {
			Description:   "The Scoped-ID of the Role Definition.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"role_definition_name"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"role_definition_name": {
			Description:   "The name of a built-in Role.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"role_definition_id"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"skip_service_principal_aad_check": {
			Description: "If the `principal_id` is a newly provisioned `Service Principal` set this value to `true` to skip the `Azure Active Directory` check which may fail due to replication lag. This argument is only valid if the `principal_id` is a `Service Principal` identity.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"delegated_managed_identity_resource_id": {
			Description:  "The delegated Azure Resource ID which contains a Managed Identity.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"description": {
			Description:  "The description for this Role Assignment.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"condition": {
			Description:  "The condition that limits the resources that the role can be assigned to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"condition_version"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"condition_version": {
			Description:  "The version of the condition.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"condition"},
			ValidateFunc: validation.StringInSlice([]string{
				"1.0",
				"2.0",
			}, false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_type": {
			Description: "The type of the `principal_id`, e.g. User, Group, Service Principal, Application, etc.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
