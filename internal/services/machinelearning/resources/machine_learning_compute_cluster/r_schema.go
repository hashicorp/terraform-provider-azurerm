// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_compute_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/machinelearningcomputes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func machineLearningComputeClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ComputeClusterName,
		},

		"machine_learning_workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"vm_size": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"vm_priority": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(machinelearningcomputes.PossibleValuesForVMPriority(), false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"scale_settings": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"max_node_count": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"min_node_count": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"scale_down_nodes_after_idle_duration": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"local_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"node_public_ip_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"ssh": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
					"admin_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Sensitive:    true,
						AtLeastOneOf: []string{"ssh.0.admin_password", "ssh.0.key_value"},
					},
					"key_value": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						AtLeastOneOf: []string{"ssh.0.admin_password", "ssh.0.key_value"},
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"ssh_public_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"subnet_resource_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C as you don't have to specify it for Azure to assign one to the cluster
			Computed: true,
			ForceNew: true,
		},

		"tags": commonschema.Tags(),
	}
}
