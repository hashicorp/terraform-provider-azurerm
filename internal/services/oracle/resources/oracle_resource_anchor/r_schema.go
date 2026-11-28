// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_resource_anchor

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Oracle Resource Anchor.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ResourceAnchorName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"linked_compartment_id": {
			Description: "Oracle Cloud Infrastructure compartment [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) which was created or linked by customer with Resource Anchor.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
