// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package capacity_reservation_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func capacityReservationGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CapacityReservationGroupName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
