// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_application_gateway_backend_address_pool_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	ipConfigurationName := d.Get("ip_configuration_name").(string)

	networkInterfaceId, err := commonids.ParseNetworkInterfaceID(d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(strings.ToLower(networkInterfaceId.ID()))
	defer locks.UnlockByID(strings.ToLower(networkInterfaceId.ID()))

	resp, err := client.Get(ctx, *networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *networkInterfaceId)
		}

		return fmt.Errorf("retrieving %s: %+v", *networkInterfaceId, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", networkInterfaceId)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", networkInterfaceId)
	}
	if resp.Model.Properties.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.ipConfigurations` was nil", networkInterfaceId)
	}
	props := resp.Model.Properties

	config := helpers.FindNetworkInterfaceIPConfiguration(resp.Model.Properties.IPConfigurations, ipConfigurationName)
	if config == nil {
		return fmt.Errorf("IP Configuration %q was not found on %s", ipConfigurationName, *networkInterfaceId)
	}
	if config.Properties == nil {
		return fmt.Errorf("`retrieving %s: ipConfiguration.properties` was nil", *networkInterfaceId)
	}
	ipConfigProps := config.Properties

	pools := make([]networkinterfaces.ApplicationGatewayBackendAddressPool, 0)

	ipConfigId := commonids.NewNetworkInterfaceIPConfigurationID(networkInterfaceId.SubscriptionId, networkInterfaceId.ResourceGroupName, networkInterfaceId.NetworkInterfaceName, ipConfigurationName)
	backendAddressPoolId, err := parse.ApplicationGatewayBackendAddressPoolID(d.Get("backend_address_pool_id").(string))
	if err != nil {
		return err
	}

	id := commonids.NewCompositeResourceID(&ipConfigId, backendAddressPoolId)

	// first double-check it doesn't exist
	exists := false
	if ipConfigProps.ApplicationGatewayBackendAddressPools != nil {
		for _, existingPool := range *ipConfigProps.ApplicationGatewayBackendAddressPools {
			if poolId := existingPool.Id; poolId != nil {
				if *poolId == backendAddressPoolId.ID() {
					exists = true
					if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError("azurerm_network_interface_application_gateway_backend_address_pool_association", id.ID())
					}
				}

				pools = append(pools, existingPool)
			}
		}
	}

	pool := networkinterfaces.ApplicationGatewayBackendAddressPool{
		Id: pointer.To(backendAddressPoolId.ID()),
	}
	if !exists {
		pools = append(pools, pool)
	}
	ipConfigProps.ApplicationGatewayBackendAddressPools = &pools

	props.IPConfigurations = helpers.UpdateNetworkInterfaceIPConfiguration(*config, props.IPConfigurations)

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, *networkInterfaceId, *resp.Model, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("updating Application Gateway Backend Address Pool Association for %s: %+v", *networkInterfaceId, err)
	}

	d.SetId(id.ID())

	return resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationRead(d, meta)
}
