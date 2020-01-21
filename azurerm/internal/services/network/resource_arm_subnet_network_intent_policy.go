package network

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSubnetNetworkIntentPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSubnetNetworkIntentPolicyCreate,
		Read:   resourceArmSubnetNetworkIntentPolicyRead,
		Delete: resourceArmSubnetNetworkIntentPolicyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmSubnetNetworkIntentPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Subnet <-> Network Security Group Association creation.")

	subnetId := d.Get("subnet_id").(string)

	parsedSubnetId, err := azure.ParseAzureResourceID(subnetId)
	if err != nil {
		return err
	}

	subnetName := parsedSubnetId.Path["subnets"]
	virtualNetworkName := parsedSubnetId.Path["virtualNetworks"]
	resourceGroup := parsedSubnetId.ResourceGroup
	serviceName := d.Get("service_name").(string)

	locks.ByName(subnetName, SubnetResourceName)
	defer locks.UnlockByName(subnetName, SubnetResourceName)

	locks.ByName(virtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, VirtualNetworkResourceName)

	req := network.PrepareNetworkPoliciesRequest{
		ServiceName: utils.String(serviceName),
	}

	future, err := client.PrepareNetworkPolicies(ctx, resourceGroup, virtualNetworkName, subnetName, req)
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

	id := fmt.Sprintf("%s/serviceName/%s", *read.ID, serviceName)

	d.SetId(id)

	return nil
}

func resourceArmSubnetNetworkIntentPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmSubnetNetworkIntentPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	virtualNetworkName := id.Path["virtualNetworks"]
	subnetName := id.Path["subnets"]
	serviceName := id.Path["serviceName"]

	locks.ByName(virtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(virtualNetworkName, VirtualNetworkResourceName)

	locks.ByName(subnetName, SubnetResourceName)
	defer locks.UnlockByName(subnetName, SubnetResourceName)

	req := network.UnprepareNetworkPoliciesRequest{
		ServiceName: utils.String(serviceName),
	}

	future, err := client.UnprepareNetworkPolicies(ctx, resourceGroup, virtualNetworkName, subnetName, req)
	if err != nil {
		return fmt.Errorf("Error removing Network Security Group Association from Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for removal of Network Security Group Association from Subnet %q (Virtual Network %q / Resource Group %q): %+v", subnetName, virtualNetworkName, resourceGroup, err)
	}

	return nil
}
