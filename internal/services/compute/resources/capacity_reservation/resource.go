// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package capacity_reservation

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name capacity_reservation -service-package-name compute -properties "name" -compare-values "subscription_id:capacity_reservation_group_id,resource_group_name:capacity_reservation_group_id,capacity_reservation_group_name:capacity_reservation_group_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCapacityReservationCreate,
		Read:   resourceCapacityReservationRead,
		Update: resourceCapacityReservationUpdate,
		Delete: resourceCapacityReservationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&capacityreservations.CapacityReservationId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&capacityreservations.CapacityReservationId{}),
		},

		Schema: capacityReservationSchema(),
	}
}
