// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (m Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"automation_account_name": {
			Description:  "The name of the Automation Account in which the Hybrid Worker is created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"worker_group_name": {
			Description:  "The name of the HybridWorker Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"worker_id": {
			Description:  "The ID of this HybridWorker in UUID notation.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"vm_resource_id": {
			Description:  "The ID of the Virtual Machine used for this HybridWorker.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"ip": {
			Description: "The IP address of assigned machine.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"registration_date_time": {
			Description: "The registration time of the worker machine.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"last_seen_date_time": {
			Description: "Last Heartbeat from the Worker.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"worker_name": {
			Description: "The name of HybridWorker.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"worker_type": {
			Description: "The type of the HybridWorker, the.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
