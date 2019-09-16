package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	networkSvc "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSubnetNetworkSecurityGroupAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSubnetNetworkSecurityGroupAssociationCreate,
		Read:   resourceArmSubnetNetworkSecurityGroupAssociationRead,
		Delete: resourceArmSubnetNetworkSecurityGroupAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"network_security_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmSubnetNetworkSecurityGroupAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.SubnetsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Subnet <-> Network Security Group Association creation.")

	subnetId := d.Get("subnet_id").(string)
	networkSecurityGroupId := d.Get("network_security_group_id").(string)

	parsedSubnetId, err := azure.ParseAzureResourceID(subnetId)
	if err != nil {
		return err
	}

	parsedNetworkSecurityGroupId, err := networkSvc.ParseNetworkSecurityGroupResourceID(networkSecurityGroupId)
	if err != nil {
		return err
	}

	locks.ByName(parsedNetworkSecurityGroupId.Name, networkSecurityGroupResourceName)
	defer locks.UnlockByName(parsedNetworkSecurityGroupId.Name, networkSecurityGroupResourceName)

	subnetName := parsedSubnetId.Path["subnets"]
	virtualNetworkName := parsedSubnetId.Path["virtualNetworks"]
	resourceGroup := parsedSubnetId.ResourceGroup

	locks.ByName(virtualNetworkName, virtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, virtualNetworkResourceName)

	locks.ByName(subnetName, subnetResourceName)
	defer locks.UnlockByName(subnetName, subnetResourceName)

	subnet, err := client.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
	if err != nil {
		if utils.ResponseWasNotFound(subnet.Response) {
			return fmt.Errorf("Subnet %q (Virtual Network %q / Resource Group %q) was not found!", subnetName, virtualNetworkName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	if props := subnet.SubnetPropertiesFormat; props != nil {
		if features.ShouldResourcesBeImported() {
			if nsg := props.NetworkSecurityGroup; nsg != nil {
				// we're intentionally not checking the ID - if there's a NSG, it needs to be imported
				if nsg.ID != nil && subnet.ID != nil {
					return tf.ImportAsExistsError("azurerm_subnet_network_security_group_association", *subnet.ID)
				}
			}
		}

		props.NetworkSecurityGroup = &network.SecurityGroup{
			ID: utils.String(networkSecurityGroupId),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualNetworkName, subnetName, subnet)
	if err != nil {
		return fmt.Errorf("Error updating Route Table Association for Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Route Table Association for Subnet %q (VN %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmSubnetNetworkSecurityGroupAssociationRead(d, meta)
}

func resourceArmSubnetNetworkSecurityGroupAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.SubnetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	virtualNetworkName := id.Path["virtualNetworks"]
	subnetName := id.Path["subnets"]

	resp, err := client.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) could not be found - removing from state!", subnetName, virtualNetworkName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	props := resp.SubnetPropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for Subnet %q (Virtual Network %q / Resource Group %q)", subnetName, virtualNetworkName, resourceGroup)
	}

	securityGroup := props.NetworkSecurityGroup
	if securityGroup == nil {
		log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) doesn't have a Network Security Group - removing from state!", subnetName, virtualNetworkName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("subnet_id", resp.ID)
	d.Set("network_security_group_id", securityGroup.ID)

	return nil
}

func resourceArmSubnetNetworkSecurityGroupAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.SubnetsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	virtualNetworkName := id.Path["virtualNetworks"]
	subnetName := id.Path["subnets"]

	// retrieve the subnet
	read, err := client.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) could not be found - removing from state!", subnetName, virtualNetworkName, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	props := read.SubnetPropertiesFormat
	if props == nil {
		return fmt.Errorf("`Properties` was nil for Subnet %q (Virtual Network %q / Resource Group %q)", subnetName, virtualNetworkName, resourceGroup)
	}

	if props.NetworkSecurityGroup == nil || props.NetworkSecurityGroup.ID == nil {
		log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) has no Network Security Group - removing from state!", subnetName, virtualNetworkName, resourceGroup)
		return nil
	}

	// once we have the network security group id to lock on, lock on that
	parsedNetworkSecurityGroupId, err := networkSvc.ParseNetworkSecurityGroupResourceID(*props.NetworkSecurityGroup.ID)
	if err != nil {
		return err
	}

	locks.ByName(parsedNetworkSecurityGroupId.Name, networkSecurityGroupResourceName)
	defer locks.UnlockByName(parsedNetworkSecurityGroupId.Name, networkSecurityGroupResourceName)

	locks.ByName(virtualNetworkName, virtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, virtualNetworkResourceName)

	locks.ByName(subnetName, subnetResourceName)
	defer locks.UnlockByName(subnetName, subnetResourceName)

	// then re-retrieve it to ensure we've got the latest state
	read, err = client.Get(ctx, resourceGroup, virtualNetworkName, subnetName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) could not be found - removing from state!", subnetName, virtualNetworkName, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	read.SubnetPropertiesFormat.NetworkSecurityGroup = nil

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualNetworkName, subnetName, read)
	if err != nil {
		return fmt.Errorf("Error removing Network Security Group Association from Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for removal of Network Security Group Association from Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	return nil
}
