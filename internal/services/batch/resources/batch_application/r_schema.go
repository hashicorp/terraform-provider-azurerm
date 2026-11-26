// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package batch_application

import (
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func batchApplicationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the application. This must be unique within the account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApplicationName,
		},

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"account_name": {
			Description:  "The name of the Batch Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"allow_updates": {
			Description: "A value indicating whether packages within the application may be overwritten using the same version string.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"default_version": {
			Description:  "The package to use if a client requests the application but does not specify a version. This property can only be set to the name of an existing package.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ApplicationVersion,
		},

		"display_name": {
			Description:  "The display name for the application.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 1024),
		},
	}
}
