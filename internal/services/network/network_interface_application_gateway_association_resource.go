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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation() *pluginsdk.Resource {
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
				ValidateFunc: validate.BackendAddressPoolID,
			},
		},
	}
}

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkInterfaces
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	ipConfigurationName := d.Get("ip_configuration_name").(string)

	networkInterfaceId, err := commonids.ParseNetworkInterfaceID(d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(networkInterfaceId.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(networkInterfaceId.NetworkInterfaceName, networkInterfaceResourceName)

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

	config := FindNetworkInterfaceIPConfiguration(resp.Model.Properties.IPConfigurations, ipConfigurationName)
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
	if ipConfigProps.ApplicationGatewayBackendAddressPools != nil {
		for _, existingPool := range *ipConfigProps.ApplicationGatewayBackendAddressPools {
			if poolId := existingPool.Id; poolId != nil {
				if *poolId == backendAddressPoolId.ID() {
					return tf.ImportAsExistsError("azurerm_network_interface_application_gateway_backend_address_pool_association", id.ID())
				}

				pools = append(pools, existingPool)
			}
		}
	}

	pool := networkinterfaces.ApplicationGatewayBackendAddressPool{
		Id: pointer.To(backendAddressPoolId.ID()),
	}
	pools = append(pools, pool)
	ipConfigProps.ApplicationGatewayBackendAddressPools = &pools

	props.IPConfigurations = updateNetworkInterfaceIPConfiguration(*config, props.IPConfigurations)

	if err := client.CreateOrUpdateThenPoll(ctx, *networkInterfaceId, *resp.Model); err != nil {
		return fmt.Errorf("updating Application Gateway Backend Address Pool Association for %s: %+v", *networkInterfaceId, err)
	}

	d.SetId(id.ID())

	return resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationRead(d, meta)
}

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceIPConfigurationId{}, &parse.ApplicationGatewayBackendAddressPoolId{})
	if err != nil {
		return err
	}

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	resp, err := client.Get(ctx, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state!", networkInterfaceId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", networkInterfaceId, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			config := FindNetworkInterfaceIPConfiguration(props.IPConfigurations, id.First.IpConfigurationName)
			if config == nil {
				log.Printf("%s was not found - removing from state!", id.First.ID())
				d.SetId("")
				return nil
			}

			found := false
			if ipConfigProps := config.Properties; ipConfigProps != nil {
				if backendPools := ipConfigProps.ApplicationGatewayBackendAddressPools; backendPools != nil {
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

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceIPConfigurationId{}, &parse.ApplicationGatewayBackendAddressPoolId{})
	if err != nil {
		return err
	}

	locks.ByName(id.First.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(id.First.NetworkInterfaceName, networkInterfaceResourceName)

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	resp, err := client.Get(ctx, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", networkInterfaceId)
		}

		return fmt.Errorf("retrieving %s : %+v", networkInterfaceId, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	if resp.Model.Properties.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.IPConfigurations` was nil", id)
	}

	config := FindNetworkInterfaceIPConfiguration(resp.Model.Properties.IPConfigurations, id.First.IpConfigurationName)
	if config == nil {
		return fmt.Errorf("IP Configuration %q was not found on %s", id.First.IpConfigurationName, id)
	}
	if config.Properties == nil {
		return fmt.Errorf("`IPConfiguration.properties` was nil for %s", id)
	}
	props := resp.Model.Properties

	ipConfigProps := config.Properties

	backendAddressPools := make([]networkinterfaces.ApplicationGatewayBackendAddressPool, 0)
	if backendPools := ipConfigProps.ApplicationGatewayBackendAddressPools; backendPools != nil {
		for _, pool := range *backendPools {
			if pool.Id == nil {
				continue
			}

			if *pool.Id != id.Second.ID() {
				backendAddressPools = append(backendAddressPools, pool)
			}
		}
	}
	ipConfigProps.ApplicationGatewayBackendAddressPools = &backendAddressPools
	props.IPConfigurations = updateNetworkInterfaceIPConfiguration(*config, props.IPConfigurations)

	if err := client.CreateOrUpdateThenPoll(ctx, networkInterfaceId, *resp.Model); err != nil {
		return fmt.Errorf("removing %s Association for %s: %+v", id.Second, id.First, err)
	}

	return nil
}
