package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationCreate,
		Read:   resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationRead,
		Delete: resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			splitId := strings.Split(id, "|")
			if _, err := parse.NetworkInterfaceIpConfigurationID(splitId[0]); err != nil {
				return err
			}

			if _, err := parse.BackendAddressPoolID(splitId[1]); err != nil {
				return err
			}
			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"network_interface_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
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
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Network Interface <-> Application Gateway Backend Address Pool Association creation.")

	networkInterfaceId := d.Get("network_interface_id").(string)
	ipConfigurationName := d.Get("ip_configuration_name").(string)
	backendAddressPoolId := d.Get("backend_address_pool_id").(string)

	id, err := parse.NetworkInterfaceID(networkInterfaceId)
	if err != nil {
		return err
	}

	locks.ByName(id.Name, networkInterfaceResourceName)
	defer locks.UnlockByName(id.Name, networkInterfaceResourceName)

	read, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("%s was not found!", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	props := read.InterfacePropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for %s", *id)
	}

	ipConfigs := props.IPConfigurations
	if ipConfigs == nil {
		return fmt.Errorf("Error: `properties.IPConfigurations` was nil for %s", *id)
	}

	c := FindNetworkInterfaceIPConfiguration(props.IPConfigurations, ipConfigurationName)
	if c == nil {
		return fmt.Errorf("Error: IP Configuration %q was not found on %s", ipConfigurationName, *id)
	}

	config := *c
	p := config.InterfaceIPConfigurationPropertiesFormat
	if p == nil {
		return fmt.Errorf("Error: `IPConfiguration.properties` was nil for %s", *id)
	}

	pools := make([]network.ApplicationGatewayBackendAddressPool, 0)

	// first double-check it doesn't exist
	resourceId := fmt.Sprintf("%s/ipConfigurations/%s|%s", networkInterfaceId, ipConfigurationName, backendAddressPoolId)
	if p.ApplicationGatewayBackendAddressPools != nil {
		for _, existingPool := range *p.ApplicationGatewayBackendAddressPools {
			if id := existingPool.ID; id != nil {
				if *id == backendAddressPoolId {
					return tf.ImportAsExistsError("azurerm_network_interface_application_gateway_backend_address_pool_association", resourceId)
				}

				pools = append(pools, existingPool)
			}
		}
	}

	pool := network.ApplicationGatewayBackendAddressPool{
		ID: utils.String(backendAddressPoolId),
	}
	pools = append(pools, pool)
	p.ApplicationGatewayBackendAddressPools = &pools

	props.IPConfigurations = updateNetworkInterfaceIPConfiguration(config, props.IPConfigurations)

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, read)
	if err != nil {
		return fmt.Errorf("updating Application Gateway Backend Address Pool Association for %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Application Gateway Backend Address Pool Association for %s: %+v", *id, err)
	}

	d.SetId(resourceId)

	return resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationRead(d, meta)
}

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}/ipConfigurations/{ipConfigurationName}|{backendAddressPoolId} but got %q", d.Id())
	}

	nicID, err := parse.NetworkInterfaceIpConfigurationID(splitId[0])
	if err != nil {
		return err
	}

	backendAddressPoolId := splitId[1]

	read, err := client.Get(ctx, nicID.ResourceGroup, nicID.NetworkInterfaceName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("%s was not found - removing from state!", *nicID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *nicID, err)
	}

	nicProps := read.InterfacePropertiesFormat
	if nicProps == nil {
		return fmt.Errorf("Error: `properties` was nil for %s", *nicID)
	}

	ipConfigs := nicProps.IPConfigurations
	if ipConfigs == nil {
		return fmt.Errorf("Error: `properties.IPConfigurations` was nil for %s", *nicID)
	}

	c := FindNetworkInterfaceIPConfiguration(nicProps.IPConfigurations, nicID.IpConfigurationName)
	if c == nil {
		log.Printf("%s was not found - removing from state!", *nicID)
		d.SetId("")
		return nil
	}
	config := *c

	found := false
	if props := config.InterfaceIPConfigurationPropertiesFormat; props != nil {
		if backendPools := props.ApplicationGatewayBackendAddressPools; backendPools != nil {
			for _, pool := range *backendPools {
				if pool.ID == nil {
					continue
				}

				if *pool.ID == backendAddressPoolId {
					found = true
					break
				}
			}
		}
	}

	if !found {
		log.Printf("[DEBUG] Association between %s and Application Gateway Backend Pool %q was not found - removing from state!", *nicID, backendAddressPoolId)
		d.SetId("")
		return nil
	}

	d.Set("backend_address_pool_id", backendAddressPoolId)
	d.Set("ip_configuration_name", nicID.IpConfigurationName)
	d.Set("network_interface_id", read.ID)

	return nil
}

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}/ipConfigurations/{ipConfigurationName}|{backendAddressPoolId} but got %q", d.Id())
	}

	nicID, err := parse.NetworkInterfaceIpConfigurationID(splitId[0])
	if err != nil {
		return err
	}

	backendAddressPoolId := splitId[1]

	locks.ByName(nicID.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(nicID.NetworkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, nicID.ResourceGroup, nicID.NetworkInterfaceName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("%s was not found!", *nicID)
		}

		return fmt.Errorf("retrieving %s : %+v", *nicID, err)
	}

	nicProps := read.InterfacePropertiesFormat
	if nicProps == nil {
		return fmt.Errorf("Error: `properties` was nil for %s", *nicID)
	}

	ipConfigs := nicProps.IPConfigurations
	if ipConfigs == nil {
		return fmt.Errorf("Error: `properties.IPConfigurations` was nil for %s", *nicID)
	}

	c := FindNetworkInterfaceIPConfiguration(nicProps.IPConfigurations, nicID.IpConfigurationName)
	if c == nil {
		return fmt.Errorf("Error: IP Configuration %q was not found on Network Interface %q (Resource Group %q)", nicID.IpConfigurationName, nicID.NetworkInterfaceName, nicID.ResourceGroup)
	}
	config := *c

	props := config.InterfaceIPConfigurationPropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: Properties for IPConfiguration %q was nil for Network Interface %q (Resource Group %q)", nicID.IpConfigurationName, nicID.NetworkInterfaceName, nicID.ResourceGroup)
	}

	backendAddressPools := make([]network.ApplicationGatewayBackendAddressPool, 0)
	if backendPools := props.ApplicationGatewayBackendAddressPools; backendPools != nil {
		for _, pool := range *backendPools {
			if pool.ID == nil {
				continue
			}

			if *pool.ID != backendAddressPoolId {
				backendAddressPools = append(backendAddressPools, pool)
			}
		}
	}
	props.ApplicationGatewayBackendAddressPools = &backendAddressPools
	nicProps.IPConfigurations = updateNetworkInterfaceIPConfiguration(config, nicProps.IPConfigurations)

	future, err := client.CreateOrUpdate(ctx, nicID.ResourceGroup, nicID.NetworkInterfaceName, read)
	if err != nil {
		return fmt.Errorf("removing Application Gateway Backend Address Pool Association for %s: %+v", *nicID, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for removal of Application Gateway Backend Address Pool Association for %s: %+v", *nicID, err)
	}

	return nil
}
