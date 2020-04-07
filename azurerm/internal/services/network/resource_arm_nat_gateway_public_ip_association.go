package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNatGatewayPublicIpAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNatGatewayPublicIpAssociationCreate,
		Read:   resourceArmNatGatewayPublicIpAssociationRead,
		Delete: resourceArmNatGatewayPublicIpAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"public_ip_address_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmNatGatewayPublicIpAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Nat Gateway <-> Public Ip Association creation.")
	natGatewayId := d.Get("nat_gateway_id").(string)
	publicIpAddressId := d.Get("public_ip_address_id").(string)
	parsedNatGatewayId, err := azure.ParseAzureResourceID(natGatewayId)
	if err != nil {
		return fmt.Errorf("Error parsing nat_gateway_id '%s': %+v", natGatewayId, err)
	}

	natGatewayName := parsedNatGatewayId.Path["natGateways"]
	resourceGroup := parsedNatGatewayId.ResourceGroup

	parsedPublicIpAddressId, err := azure.ParseAzureResourceID(publicIpAddressId)
	if err != nil {
		return fmt.Errorf("Error parsing public_ip_address_id '%s': %+v", publicIpAddressId, err)
	}

	publicIpAddressName := parsedPublicIpAddressId.Path["publicIPAddresses"]

	locks.ByName(natGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(natGatewayName, natGatewayResourceName)
	locks.ByName(publicIpAddressName, publicIpResourceName)
	defer locks.UnlockByName(publicIpAddressName, publicIpResourceName)

	natGateway, err := client.Get(ctx, resourceGroup, natGatewayName, "")
	if err != nil {
		if utils.ResponseWasNotFound(natGateway.Response) {
			return fmt.Errorf("Nat Gateway %q (Resource Group %q) was not found!", natGatewayName, resourceGroup)
		}
		return fmt.Errorf("Error retrieving Nat Gateway %q (Resource Group %q): %+v", natGatewayName, resourceGroup, err)
	}

	publicIpAddresses := make([]network.SubResource, 0)

	if natGateway.PublicIPAddresses != nil {
		for _, existingPublicIPAddress := range *natGateway.PublicIPAddresses {
			if id := existingPublicIPAddress.ID; id != nil {
				if *id == publicIpAddressId {
					if features.ShouldResourcesBeImported() {
						return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_association", *natGateway.ID)
					}

					continue
				}

				publicIpAddresses = append(publicIpAddresses, existingPublicIPAddress)
			}
		}
	}

	publicIpAddress := network.SubResource{
		ID: utils.String(publicIpAddressId),
	}
	publicIpAddresses = append(publicIpAddresses, publicIpAddress)
	natGateway.PublicIPAddresses = &publicIpAddresses

	future, err := client.CreateOrUpdate(ctx, resourceGroup, natGatewayName, natGateway)
	if err != nil {
		return fmt.Errorf("Error updating Public IP Association for Nat Gateway %q (Resource Group %q): %+v", natGatewayName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Public IP Association for Nat Gateway %q (Resource Group %q): %+v", natGatewayName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, natGatewayName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Nat Gateway %q (Resource Group %q): %+v", natGatewayName, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmNatGatewayPublicIpAssociationRead(d, meta)
}

func resourceArmNatGatewayPublicIpAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	natGatewayName := id.Path["natGateways"]
	publicIpAddressId := d.Get("public_ip_address_id").(string)

	natGateway, err := client.Get(ctx, resourceGroup, natGatewayName, "")
	if err != nil {
		if utils.ResponseWasNotFound(natGateway.Response) {
			log.Printf("[DEBUG] Nat Gateway %q (Resource Group %q) could not be found - removing from state!", natGatewayName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Nat Gateway %q (Resource Group %q): %+v", natGatewayName, resourceGroup, err)
	}

	props := natGateway.NatGatewayPropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for Nat Gateway %q (Resource Group %q)", natGatewayName, resourceGroup)
	}
	publicIpAddresses := props.PublicIPAddresses
	if publicIpAddresses == nil {
		log.Printf("[DEBUG] Nat Gateway %q (Resource Group %q) doesn't have a Public IP - removing from state!", natGatewayName, resourceGroup)
		d.SetId("")
		return nil
	}

	found := false
	if props := natGateway.NatGatewayPropertiesFormat; props != nil {
		if publicIPAddresses := props.PublicIPAddresses; publicIPAddresses != nil {
			for _, publicIPAddress := range *publicIPAddresses {
				if publicIPAddress.ID == nil {
					continue
				}

				if *publicIPAddress.ID == publicIpAddressId {
					found = true
					break
				}
			}
		}
	}

	if !found {
		log.Printf("[DEBUG] Association between Nat Gateway %q (Resource Group %q) and Public IP %q was not found - removing from state!", natGatewayName, resourceGroup, publicIpAddressId)
		d.SetId("")
		return nil
	}

	d.Set("nat_gateway_id", natGateway.ID)
	d.Set("public_ip_address_id", publicIpAddressId)

	return nil
}

func resourceArmNatGatewayPublicIpAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	natGatewayName := id.Path["natGateways"]
	publicIpAddressId := d.Get("public_ip_address_id").(string)

	locks.ByName(natGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(natGatewayName, natGatewayResourceName)

	natGateway, err := client.Get(ctx, resourceGroup, natGatewayName, "")
	if err != nil {
		if utils.ResponseWasNotFound(natGateway.Response) {
			return fmt.Errorf("Nat Gateway %q (Resource Group %q) was not found!", natGatewayName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Nat Gateway %q (Resource Group %q): %+v", natGatewayName, resourceGroup, err)
	}

	props := natGateway.NatGatewayPropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for Nat Gateway %q (Resource Group %q)", natGatewayName, resourceGroup)
	}

	publicIpAddresses := make([]network.SubResource, 0)
	if publicIPAddresses := props.PublicIPAddresses; publicIPAddresses != nil {
		for _, publicIPAddress := range *publicIPAddresses {
			if publicIPAddress.ID == nil {
				continue
			}

			if *publicIPAddress.ID != publicIpAddressId {
				publicIpAddresses = append(publicIpAddresses, publicIPAddress)
			}
		}
	}
	props.PublicIPAddresses = &publicIpAddresses

	future, err := client.CreateOrUpdate(ctx, resourceGroup, natGatewayName, natGateway)
	if err != nil {
		return fmt.Errorf("Error removing Public Ip Association for Nat Gateway %q (Resource Group %q): %+v", natGatewayName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for removal of Public Ip Association for Nat Gateway %q (Resource Group %q): %+v", natGatewayName, resourceGroup, err)
	}

	return nil
}
