// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceBackendAddressPoolAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceBackendAddressPoolAssociationCreate,
		Read:   resourceNetworkInterfaceBackendAddressPoolAssociationRead,
		Delete: resourceNetworkInterfaceBackendAddressPoolAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &commonids.NetworkInterfaceIPConfigurationId{}, &loadbalancers.BackendAddressPoolId{})
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"network_interface_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateNetworkInterfaceID,
			},

			"ip_configuration_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"backend_address_pool_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loadbalancers.ValidateBackendAddressPoolID,
			},
		},
	}
}

func resourceNetworkInterfaceBackendAddressPoolAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkInterfaces
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

	locks.ByName(networkInterfaceId.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(networkInterfaceId.NetworkInterfaceName, networkInterfaceResourceName)

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

	config := FindNetworkInterfaceIPConfiguration(read.Model.Properties.IPConfigurations, ipConfigId.IpConfigurationName)
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
	if ipConfigProps.LoadBalancerBackendAddressPools != nil {
		for _, existingPool := range *ipConfigProps.LoadBalancerBackendAddressPools {
			if poolId := existingPool.Id; poolId != nil {
				if *poolId == backendAddressPoolId.ID() {
					return tf.ImportAsExistsError("azurerm_network_interface_backend_address_pool_association", id.ID())
				}

				pools = append(pools, existingPool)
			}
		}
	}

	pool := networkinterfaces.BackendAddressPool{
		Id: pointer.To(backendAddressPoolId.ID()),
	}
	pools = append(pools, pool)
	ipConfigProps.LoadBalancerBackendAddressPools = &pools

	props.IPConfigurations = updateNetworkInterfaceIPConfiguration(*config, props.IPConfigurations)

	if err := client.CreateOrUpdateThenPoll(ctx, *networkInterfaceId, *read.Model); err != nil {
		return fmt.Errorf("updating Backend Address Pool Association for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkInterfaceBackendAddressPoolAssociationRead(d, meta)
}

func resourceNetworkInterfaceBackendAddressPoolAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceIPConfigurationId{}, &loadbalancers.BackendAddressPoolId{})
	if err != nil {
		return err
	}

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	read, err := client.Get(ctx, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("%s was not found - removing from state!", networkInterfaceId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", networkInterfaceId, err)
	}

	if model := read.Model; model != nil {
		if props := model.Properties; props != nil {
			config := FindNetworkInterfaceIPConfiguration(props.IPConfigurations, id.First.IpConfigurationName)
			if config == nil {
				log.Printf("IP Configuration %q was not found in %s - removing from state!", id.First.IpConfigurationName, networkInterfaceId)
				d.SetId("")
				return nil
			}

			found := false
			if ipConfigProps := config.Properties; ipConfigProps != nil {
				if backendPools := ipConfigProps.LoadBalancerBackendAddressPools; backendPools != nil {
					for _, pool := range *backendPools {
						if pool.Id == nil {
							continue
						}

						if *pool.Id == id.Second.ID() {
							found = true
							break
						}
					}
				}
			}
			if !found {
				log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.First, id.Second)
				d.SetId("")
				return nil
			}
		}
	}

	d.Set("backend_address_pool_id", id.Second.ID())
	d.Set("ip_configuration_name", id.First.IpConfigurationName)
	d.Set("network_interface_id", networkInterfaceId.ID())

	return nil
}

func resourceNetworkInterfaceBackendAddressPoolAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceIPConfigurationId{}, &loadbalancers.BackendAddressPoolId{})
	if err != nil {
		return err
	}

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	locks.ByName(id.First.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(id.First.NetworkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
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
	config := FindNetworkInterfaceIPConfiguration(read.Model.Properties.IPConfigurations, id.First.IpConfigurationName)
	if config == nil {
		return fmt.Errorf("IP Configuration %q was not found on %s", id.First.IpConfigurationName, networkInterfaceId)
	}

	ipConfigProps := config.Properties
	if ipConfigProps == nil {
		return fmt.Errorf("retrieving %s: `properties` for IPConfiguration %q was nil", id.First.IpConfigurationName, networkInterfaceId)
	}

	backendAddressPools := make([]networkinterfaces.BackendAddressPool, 0)
	if backendPools := ipConfigProps.LoadBalancerBackendAddressPools; backendPools != nil {
		for _, pool := range *backendPools {
			if pool.Id == nil {
				continue
			}

			if *pool.Id != id.Second.ID() {
				backendAddressPools = append(backendAddressPools, pool)
			}
		}
	}
	ipConfigProps.LoadBalancerBackendAddressPools = &backendAddressPools
	props.IPConfigurations = updateNetworkInterfaceIPConfiguration(*config, props.IPConfigurations)

	if err := client.CreateOrUpdateThenPoll(ctx, networkInterfaceId, *read.Model); err != nil {
		return fmt.Errorf("removing Backend Address Pool Association for %s: %+v", networkInterfaceId, err)
	}

	return nil
}
