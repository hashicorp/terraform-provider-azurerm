// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_slot_virtual_network_swift_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func appServiceSlotVirtualNetworkSwiftConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"app_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateAppServiceID,
		},
		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},
		"slot_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AppServiceName,
		},
	}
}
