package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetworkInterfaceApplicationSecurityGroupAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkInterfaceApplicationSecurityGroupAssociationCreate,
		Read:   resourceArmNetworkInterfaceApplicationSecurityGroupAssociationRead,
		Delete: resourceArmNetworkInterfaceApplicationSecurityGroupAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"ip_configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"application_security_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmNetworkInterfaceApplicationSecurityGroupAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.InterfacesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Network Interface <-> Application Security Group Association creation.")

	networkInterfaceId := d.Get("network_interface_id").(string)
	ipConfigurationName := d.Get("ip_configuration_name").(string)
	applicationSecurityGroupId := d.Get("application_security_group_id").(string)

	id, err := azure.ParseAzureResourceID(networkInterfaceId)
	if err != nil {
		return err
	}

	networkInterfaceName := id.Path["networkInterfaces"]
	resourceGroup := id.ResourceGroup

	locks.ByName(networkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(networkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, resourceGroup, networkInterfaceName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("Network Interface %q (Resource Group %q) was not found!", networkInterfaceName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", networkInterfaceName, resourceGroup, err)
	}

	props := read.InterfacePropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for Network Interface %q (Resource Group %q)", networkInterfaceName, resourceGroup)
	}

	ipConfigs := props.IPConfigurations
	if ipConfigs == nil {
		return fmt.Errorf("Error: `properties.IPConfigurations` was nil for Network Interface %q (Resource Group %q)", networkInterfaceName, resourceGroup)
	}

	c := azure.FindNetworkInterfaceIPConfiguration(props.IPConfigurations, ipConfigurationName)
	if c == nil {
		return fmt.Errorf("Error: IP Configuration %q was not found on Network Interface %q (Resource Group %q)", ipConfigurationName, networkInterfaceName, resourceGroup)
	}

	config := *c
	p := config.InterfaceIPConfigurationPropertiesFormat
	if p == nil {
		return fmt.Errorf("Error: `IPConfiguration.properties` was nil for Network Interface %q (Resource Group %q)", networkInterfaceName, resourceGroup)
	}

	applicationSecurityGroups := make([]network.ApplicationSecurityGroup, 0)

	// first double-check it doesn't exist
	if p.ApplicationSecurityGroups != nil {
		for _, existingGroup := range *p.ApplicationSecurityGroups {
			if id := existingGroup.ID; id != nil {
				if *id == applicationSecurityGroupId {
					if features.ShouldResourcesBeImported() {
						return tf.ImportAsExistsError("azurerm_network_interface_application_security_group_association", *id)
					}

					continue
				}

				applicationSecurityGroups = append(applicationSecurityGroups, existingGroup)
			}
		}
	}

	group := network.ApplicationSecurityGroup{
		ID: utils.String(applicationSecurityGroupId),
	}
	applicationSecurityGroups = append(applicationSecurityGroups, group)
	p.ApplicationSecurityGroups = &applicationSecurityGroups

	props.IPConfigurations = azure.UpdateNetworkInterfaceIPConfiguration(config, props.IPConfigurations)

	future, err := client.CreateOrUpdate(ctx, resourceGroup, networkInterfaceName, read)
	if err != nil {
		return fmt.Errorf("Error updating Application Security Group Association for Network Interface %q (Resource Group %q): %+v", networkInterfaceName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Application Security Group Association for NIC %q (Resource Group %q): %+v", networkInterfaceName, resourceGroup, err)
	}

	resourceId := fmt.Sprintf("%s/ipConfigurations/%s|%s", networkInterfaceId, ipConfigurationName, applicationSecurityGroupId)
	d.SetId(resourceId)

	return resourceArmNetworkInterfaceApplicationSecurityGroupAssociationRead(d, meta)
}

func resourceArmNetworkInterfaceApplicationSecurityGroupAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.InterfacesClient
	ctx := meta.(*ArmClient).StopContext

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}/ipConfigurations/{ipConfigurationName}|{applicationSecurityGroupId} but got %q", d.Id())
	}

	nicID, err := azure.ParseAzureResourceID(splitId[0])
	if err != nil {
		return err
	}

	ipConfigurationName := nicID.Path["ipConfigurations"]
	networkInterfaceName := nicID.Path["networkInterfaces"]
	resourceGroup := nicID.ResourceGroup
	applicationSecurityGroupId := splitId[1]

	read, err := client.Get(ctx, resourceGroup, networkInterfaceName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("Network Interface %q (Resource Group %q) was not found!", networkInterfaceName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", networkInterfaceName, resourceGroup, err)
	}

	nicProps := read.InterfacePropertiesFormat
	if nicProps == nil {
		return fmt.Errorf("Error: `properties` was nil for Network Interface %q (Resource Group %q)", networkInterfaceName, resourceGroup)
	}

	ipConfigs := nicProps.IPConfigurations
	if ipConfigs == nil {
		return fmt.Errorf("Error: `properties.IPConfigurations` was nil for Network Interface %q (Resource Group %q)", networkInterfaceName, resourceGroup)
	}

	c := azure.FindNetworkInterfaceIPConfiguration(nicProps.IPConfigurations, ipConfigurationName)
	if c == nil {
		log.Printf("IP Configuration %q was not found in Network Interface %q (Resource Group %q) - removing from state!", ipConfigurationName, networkInterfaceName, resourceGroup)
		d.SetId("")
		return nil
	}
	config := *c

	found := false
	if props := config.InterfaceIPConfigurationPropertiesFormat; props != nil {
		if groups := props.ApplicationSecurityGroups; groups != nil {
			for _, group := range *groups {
				if group.ID == nil {
					continue
				}

				if *group.ID == applicationSecurityGroupId {
					found = true
					break
				}
			}
		}
	}

	if !found {
		log.Printf("[DEBUG] Association between Network Interface %q (Resource Group %q) and Application Security Group %q was not found - removing from state!", networkInterfaceName, resourceGroup, applicationSecurityGroupId)
		d.SetId("")
		return nil
	}

	d.Set("application_security_group_id", applicationSecurityGroupId)
	d.Set("ip_configuration_name", ipConfigurationName)
	if id := read.ID; id != nil {
		d.Set("network_interface_id", *id)
	}

	return nil
}

func resourceArmNetworkInterfaceApplicationSecurityGroupAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.InterfacesClient
	ctx := meta.(*ArmClient).StopContext

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}/ipConfigurations/{ipConfigurationName}|{applicationSecurityGroupId} but got %q", d.Id())
	}

	nicID, err := azure.ParseAzureResourceID(splitId[0])
	if err != nil {
		return err
	}

	ipConfigurationName := nicID.Path["ipConfigurations"]
	networkInterfaceName := nicID.Path["networkInterfaces"]
	resourceGroup := nicID.ResourceGroup
	applicationSecurityGroupId := splitId[1]

	locks.ByName(networkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(networkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, resourceGroup, networkInterfaceName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("Network Interface %q (Resource Group %q) was not found!", networkInterfaceName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", networkInterfaceName, resourceGroup, err)
	}

	nicProps := read.InterfacePropertiesFormat
	if nicProps == nil {
		return fmt.Errorf("Error: `properties` was nil for Network Interface %q (Resource Group %q)", networkInterfaceName, resourceGroup)
	}

	ipConfigs := nicProps.IPConfigurations
	if ipConfigs == nil {
		return fmt.Errorf("Error: `properties.IPConfigurations` was nil for Network Interface %q (Resource Group %q)", networkInterfaceName, resourceGroup)
	}

	c := azure.FindNetworkInterfaceIPConfiguration(nicProps.IPConfigurations, ipConfigurationName)
	if c == nil {
		return fmt.Errorf("Error: IP Configuration %q was not found on Network Interface %q (Resource Group %q)", ipConfigurationName, networkInterfaceName, resourceGroup)
	}
	config := *c

	props := config.InterfaceIPConfigurationPropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: Properties for IPConfiguration %q was nil for Network Interface %q (Resource Group %q)", ipConfigurationName, networkInterfaceName, resourceGroup)
	}

	applicationSecurityGroups := make([]network.ApplicationSecurityGroup, 0)
	if groups := props.ApplicationSecurityGroups; groups != nil {
		for _, pool := range *groups {
			if pool.ID == nil {
				continue
			}

			if *pool.ID != applicationSecurityGroupId {
				applicationSecurityGroups = append(applicationSecurityGroups, pool)
			}
		}
	}
	props.ApplicationSecurityGroups = &applicationSecurityGroups
	nicProps.IPConfigurations = azure.UpdateNetworkInterfaceIPConfiguration(config, nicProps.IPConfigurations)

	future, err := client.CreateOrUpdate(ctx, resourceGroup, networkInterfaceName, read)
	if err != nil {
		return fmt.Errorf("Error removing Application Security Group for Network Interface %q (Resource Group %q): %+v", networkInterfaceName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for removal of Application Security Group for NIC %q (Resource Group %q): %+v", networkInterfaceName, resourceGroup, err)
	}

	return nil
}
