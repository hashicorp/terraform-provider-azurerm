// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_watcher

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2020-01-13-preview/watcher"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (m Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"automation_account_id": {
			Description:  "The ID of Automation Account to manage this Watcher.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: watcher.ValidateAutomationAccountID,
		},

		"name": {
			Description:  "The name which should be used for this Automation Watcher.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),

		"etag": {
			Description:  "A string of etag assigned to this Automation Watcher.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"execution_frequency_in_seconds": {
			Description:  "The frequency at which the watcher is invoked.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"script_name": {
			Description:  "The name of an existing runbook this watcher is attached to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"script_parameters": {
			Description: "A list of key-value parameters.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"script_run_on": {
			Description:  "The name of the Hybrid work group the watcher will run on.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Description:  "A description of this Automation Watcher.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"status": {
			Description: "The current status of the Automation Watcher.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
