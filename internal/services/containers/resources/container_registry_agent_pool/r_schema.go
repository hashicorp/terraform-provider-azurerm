// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_agent_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func containerRegistryAgentPoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Azure Container Registry Agent Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 20),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"container_registry_name": {
			Description:  "Name of Azure Container Registry to create an Agent Pool for.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"instance_count": {
			Description: "VMSS instance count.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     1,
		},

		"tier": {
			Description: "Sets the VM your agent pool will run on.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "S1",
			ValidateFunc: validation.StringInSlice([]string{
				"S1",
				"S2",
				"S3",
				"I6",
			}, false),
		},

		"virtual_network_subnet_id": {
			Description:  "The ID of the Virtual Network Subnet Resource where the agent machines will be running.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"tags": commonschema.Tags(),
	}
}
