// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_application_security_group_association

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationsecuritygroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/resources/network_interface_application_security_group_association/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceApplicationSecurityGroupAssociationCreate,
		Read:   resourceNetworkInterfaceApplicationSecurityGroupAssociationRead,
		Delete: resourceNetworkInterfaceApplicationSecurityGroupAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &commonids.NetworkInterfaceId{}, &applicationsecuritygroups.ApplicationSecurityGroupId{})
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NetworkInterfaceApplicationSecurityGroupAssociationV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: networkInterfaceApplicationSecurityGroupAssociationSchema(),
	}
}
