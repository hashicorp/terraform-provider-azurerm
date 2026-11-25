// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dedicated_host_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dedicatedHostGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Dedicated Host Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DedicatedHostGroupName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"platform_fault_domain_count": {
			Description:  "The number of fault domains that the Dedicated Host Group spans.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 3),
		},

		"automatic_placement_enabled": {
			Description: "Whether Virtual Machines or Virtual Machine Scale Sets are placed automatically on this Dedicated Host Group.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
		"zone": commonschema.ZoneSingleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
