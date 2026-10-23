// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_application_gateway_backend_address_pool_association

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationCreate,
		Read:   resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationRead,
		Delete: resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &commonids.NetworkInterfaceIPConfigurationId{}, &parse.ApplicationGatewayBackendAddressPoolId{})
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: networkInterfaceApplicationGatewayBackendAddressPoolAssociationSchema(),
	}
}
