// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_credential

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func automationCredentialSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "Specifies the name of the Credential.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Description:  "The name of the automation account in which the Credential is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccount(),
		},

		"username": {
			Description: "The username associated with this Automation Credential.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"password": {
			Description: "The password associated with this Automation Credential.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			Sensitive:   true,
		},

		"description": {
			Description: "The description associated with this Automation Credential.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
	}
}
