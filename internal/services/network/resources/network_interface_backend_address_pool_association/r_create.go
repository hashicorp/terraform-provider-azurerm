// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_backend_address_pool_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceBackendAddressPoolAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	networkInterfaceId, err := commonids.ParseNetworkInterfaceID(d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}
	backendAddressPoolId, err := loadbalancers.ParseBackendAddressPoolID(d.Get("backend_address_pool_id").(string))
	if err != nil {
		return err
	}
	ipConfigId := commonids.NewNetworkInterfaceIPConfigurationID(networkInterfaceId.SubscriptionId, networkInterfaceId.ResourceGroupName, networkInterfaceId.NetworkInterfaceName, d.Get("ip_configuration_name").(string))

	locks.ByID(strings.ToLower(networkInterfaceId.ID()))
	defer locks.UnlockByID(strings.ToLower(networkInterfaceId.ID()))

	read, err := client.Get(ctx, *networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf("%s was not found", networkInterfaceId)
		}
		return fmt.Errorf("retrieving %s: %+v", networkInterfaceId, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", networkInterfaceId)
	}
	if read.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", networkInterfaceId)
	}
	if read.Model.Properties.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.ipConfigurations` was nil", networkInterfaceId)
	}

	props := read.Model.Properties

	config := helpers.FindNetworkInterfaceIPConfiguration(read.Model.Properties.IPConfigurations, ipConfigId.IpConfigurationName)
	if config == nil {
		return fmt.Errorf("IP Configuration %q was not found on %s", ipConfigId.IpConfigurationName, networkInterfaceId)
	}

	ipConfigProps := config.Properties
	if ipConfigProps == nil {
		return fmt.Errorf("retrieving %s: `ipConfiguration.properties` was nil", networkInterfaceId)
	}

	id := commonids.NewCompositeResourceID(&ipConfigId, backendAddressPoolId)
	pools := make([]networkinterfaces.BackendAddressPool, 0)

	// first double-check it doesn't exist
	exists := false
	if ipConfigProps.LoadBalancerBackendAddressPools != nil {
		for _, existingPool := range *ipConfigProps.LoadBalancerBackendAddressPools {
			if poolId := existingPool.Id; poolId != nil {
				if *poolId == backendAddressPoolId.ID() {
					exists = true
					if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError("azurerm_network_interface_backend_address_pool_association", id.ID())
					}
				}

				pools = append(pools, existingPool)
			}
		}
	}

	pool := networkinterfaces.BackendAddressPool{
		Id: pointer.To(backendAddressPoolId.ID()),
	}
	if !exists {
		pools = append(pools, pool)
	}
	ipConfigProps.LoadBalancerBackendAddressPools = &pools

	props.IPConfigurations = helpers.UpdateNetworkInterfaceIPConfiguration(*config, props.IPConfigurations)

	if err := client.CreateOrUpdateThenPoll(ctx, *networkInterfaceId, *read.Model); err != nil {
		return fmt.Errorf("updating Backend Address Pool Association for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkInterfaceBackendAddressPoolAssociationRead(d, meta)
}
