// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	loadBalancerParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	loadBalancerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceNetworkInterfaceNatRuleAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceNatRuleAssociationCreate,
		Read:   resourceNetworkInterfaceNatRuleAssociationRead,
		Delete: resourceNetworkInterfaceNatRuleAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			splitId := strings.Split(id, "|")
			if _, err := parse.NetworkInterfaceIpConfigurationID(splitId[0]); err != nil {
				return err
			}
			if _, err := loadBalancerParse.LoadBalancerInboundNatRuleID(splitId[1]); err != nil {
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
				ValidateFunc: validate.NetworkInterfaceID,
			},

			"ip_configuration_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"nat_rule_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loadBalancerValidate.LoadBalancerInboundNatRuleID,
			},
		},
	}
}

func resourceNetworkInterfaceNatRuleAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Network Interface <-> Load Balancer NAT Rule Association creation.")

	networkInterfaceId := d.Get("network_interface_id").(string)
	ipConfigurationName := d.Get("ip_configuration_name").(string)
	natRuleId := d.Get("nat_rule_id").(string)

	id, err := parse.NetworkInterfaceID(networkInterfaceId)
	if err != nil {
		return err
	}

	locks.ByName(id.Name, networkInterfaceResourceName)
	defer locks.UnlockByName(id.Name, networkInterfaceResourceName)

	read, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf(" %s was not found!", *id)
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

	rules := make([]network.InboundNatRule, 0)

	// first double-check it doesn't exist
	resourceId := fmt.Sprintf("%s/ipConfigurations/%s|%s", networkInterfaceId, ipConfigurationName, natRuleId)
	if p.LoadBalancerInboundNatRules != nil {
		for _, existingRule := range *p.LoadBalancerInboundNatRules {
			if id := existingRule.ID; id != nil {
				if *id == natRuleId {
					return tf.ImportAsExistsError("azurerm_network_interface_nat_rule_association", resourceId)
				}

				rules = append(rules, existingRule)
			}
		}
	}

	rule := network.InboundNatRule{
		ID: utils.String(natRuleId),
	}
	rules = append(rules, rule)
	p.LoadBalancerInboundNatRules = &rules

	props.IPConfigurations = updateNetworkInterfaceIPConfiguration(config, props.IPConfigurations)

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, read)
	if err != nil {
		return fmt.Errorf("updating NAT Rule Association for %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of NAT Rule Association for %s: %+v", *id, err)
	}

	d.SetId(resourceId)

	return resourceNetworkInterfaceNatRuleAssociationRead(d, meta)
}

func resourceNetworkInterfaceNatRuleAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}/ipConfigurations/{ipConfigurationName}|{natRuleId} but got %q", d.Id())
	}

	nicID, err := parse.NetworkInterfaceIpConfigurationID(splitId[0])
	if err != nil {
		return err
	}

	natRuleId := splitId[1]

	read, err := client.Get(ctx, nicID.ResourceGroup, nicID.NetworkInterfaceName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("Network Interface %q (Resource Group %q) was not found - removing from state!", nicID.NetworkInterfaceName, nicID.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Network Interface %q (Resource Group %q): %+v", nicID.NetworkInterfaceName, nicID.ResourceGroup, err)
	}

	nicProps := read.InterfacePropertiesFormat
	if nicProps == nil {
		return fmt.Errorf("Error: `properties` was nil for Network Interface %q (Resource Group %q)", nicID.NetworkInterfaceName, nicID.ResourceGroup)
	}

	ipConfigs := nicProps.IPConfigurations
	if ipConfigs == nil {
		return fmt.Errorf("Error: `properties.IPConfigurations` was nil for Network Interface %q (Resource Group %q)", nicID.NetworkInterfaceName, nicID.ResourceGroup)
	}

	c := FindNetworkInterfaceIPConfiguration(nicProps.IPConfigurations, nicID.IpConfigurationName)
	if c == nil {
		log.Printf("IP Configuration %q was not found in Network Interface %q (Resource Group %q) - removing from state!", nicID.IpConfigurationName, nicID.NetworkInterfaceName, nicID.ResourceGroup)
		d.SetId("")
		return nil
	}
	config := *c

	found := false
	if props := config.InterfaceIPConfigurationPropertiesFormat; props != nil {
		if rules := props.LoadBalancerInboundNatRules; rules != nil {
			for _, rule := range *rules {
				if rule.ID == nil {
					continue
				}

				if *rule.ID == natRuleId {
					found = true
					break
				}
			}
		}
	}

	if !found {
		log.Printf("[DEBUG] Association between Network Interface %q (Resource Group %q) and Load Balancer NAT Rule %q was not found - removing from state!", nicID.NetworkInterfaceName, nicID.ResourceGroup, natRuleId)
		d.SetId("")
		return nil
	}

	d.Set("ip_configuration_name", nicID.IpConfigurationName)
	d.Set("nat_rule_id", natRuleId)
	d.Set("network_interface_id", read.ID)

	return nil
}

func resourceNetworkInterfaceNatRuleAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}/ipConfigurations/{ipConfigurationName}|{natRuleId} but got %q", d.Id())
	}

	nicID, err := parse.NetworkInterfaceIpConfigurationID(splitId[0])
	if err != nil {
		return err
	}

	natRuleId := splitId[1]

	locks.ByName(nicID.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(nicID.NetworkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, nicID.ResourceGroup, nicID.NetworkInterfaceName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("Network Interface %q (Resource Group %q) was not found!", nicID.NetworkInterfaceName, nicID.ResourceGroup)
		}

		return fmt.Errorf("retrieving Network Interface %q (Resource Group %q): %+v", nicID.NetworkInterfaceName, nicID.ResourceGroup, err)
	}

	nicProps := read.InterfacePropertiesFormat
	if nicProps == nil {
		return fmt.Errorf("Error: `properties` was nil for Network Interface %q (Resource Group %q)", nicID.NetworkInterfaceName, nicID.ResourceGroup)
	}

	ipConfigs := nicProps.IPConfigurations
	if ipConfigs == nil {
		return fmt.Errorf("Error: `properties.IPConfigurations` was nil for Network Interface %q (Resource Group %q)", nicID.NetworkInterfaceName, nicID.ResourceGroup)
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

	updatedRules := make([]network.InboundNatRule, 0)
	if existingRules := props.LoadBalancerInboundNatRules; existingRules != nil {
		for _, rule := range *existingRules {
			if rule.ID == nil {
				continue
			}

			if *rule.ID != natRuleId {
				updatedRules = append(updatedRules, rule)
			}
		}
	}
	props.LoadBalancerInboundNatRules = &updatedRules
	nicProps.IPConfigurations = updateNetworkInterfaceIPConfiguration(config, nicProps.IPConfigurations)

	future, err := client.CreateOrUpdate(ctx, nicID.ResourceGroup, nicID.NetworkInterfaceName, read)
	if err != nil {
		return fmt.Errorf("removing NAT Rule Association for Network Interface %q (Resource Group %q): %+v", nicID.NetworkInterfaceName, nicID.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for removal of NAT Rule Association for NIC %q (Resource Group %q): %+v", nicID.NetworkInterfaceName, nicID.ResourceGroup, err)
	}

	return nil
}
