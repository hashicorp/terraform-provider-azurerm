package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNatGatewayPublicIpAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNatGatewayPublicIpAssociationCreate,
		Read:   resourceArmNatGatewayPublicIpAssociationRead,
		Delete: resourceArmNatGatewayPublicIpAssociationDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.NatGatewayID(id)
			return err
		}),

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
				ValidateFunc: validate.NatGatewayID,
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
	parsedNatGatewayId, err := parse.NatGatewayID(natGatewayId)
	if err != nil {
		return err
	}

	locks.ByName(parsedNatGatewayId.Name, natGatewayResourceName)
	defer locks.UnlockByName(parsedNatGatewayId.Name, natGatewayResourceName)

	natGateway, err := client.Get(ctx, parsedNatGatewayId.ResourceGroup, parsedNatGatewayId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(natGateway.Response) {
			return fmt.Errorf("Nat Gateway %q (Resource Group %q) was not found.", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup)
		}
		return fmt.Errorf("failed to retrieve Nat Gateway %q (Resource Group %q): %+v", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup, err)
	}

	publicIpAddresses := make([]network.SubResource, 0)
	if natGateway.PublicIPAddresses != nil {
		for _, existingPublicIPAddress := range *natGateway.PublicIPAddresses {
			if existingPublicIPAddress.ID != nil {
				if *existingPublicIPAddress.ID == publicIpAddressId {
					return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_association", *natGateway.ID)
				}

				publicIpAddresses = append(publicIpAddresses, existingPublicIPAddress)
			}
		}
	}

	publicIpAddresses = append(publicIpAddresses, network.SubResource{
		ID: utils.String(publicIpAddressId),
	})
	natGateway.PublicIPAddresses = &publicIpAddresses

	future, err := client.CreateOrUpdate(ctx, parsedNatGatewayId.ResourceGroup, parsedNatGatewayId.Name, natGateway)
	if err != nil {
		return fmt.Errorf("failed to update Public IP Association for Nat Gateway %q (Resource Group %q): %+v", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed to wait for completion of Public IP Association for Nat Gateway %q (Resource Group %q): %+v", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, parsedNatGatewayId.ResourceGroup, parsedNatGatewayId.Name, "")
	if err != nil {
		return fmt.Errorf("failed to retrieve Nat Gateway %q (Resource Group %q): %+v", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmNatGatewayPublicIpAssociationRead(d, meta)
}

func resourceArmNatGatewayPublicIpAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayID(d.Id())
	if err != nil {
		return err
	}

	natGateway, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(natGateway.Response) {
			log.Printf("[DEBUG] Nat Gateway %q (Resource Group %q) could not be found - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to retrieve Nat Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if natGateway.PublicIPAddresses == nil {
		log.Printf("[DEBUG] Nat Gateway %q (Resource Group %q) doesn't have a Public IP - removing from state!", id.Name, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("nat_gateway_id", natGateway.ID)

	return nil
}

func resourceArmNatGatewayPublicIpAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayID(d.Id())
	if err != nil {
		return err
	}

	publicIpAddressId := d.Get("public_ip_address_id").(string)

	locks.ByName(id.Name, natGatewayResourceName)
	defer locks.UnlockByName(id.Name, natGatewayResourceName)

	natGateway, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(natGateway.Response) {
			return fmt.Errorf("Nat Gateway %q (Resource Group %q) was not found.", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("failed to retrieve Nat Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	publicIpAddresses := make([]network.SubResource, 0)
	if publicIPAddresses := natGateway.PublicIPAddresses; publicIPAddresses != nil {
		for _, publicIPAddress := range *publicIPAddresses {
			if publicIPAddress.ID == nil {
				continue
			}

			if *publicIPAddress.ID != publicIpAddressId {
				publicIpAddresses = append(publicIpAddresses, publicIPAddress)
			}
		}
	}
	natGateway.PublicIPAddresses = &publicIpAddresses

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, natGateway)
	if err != nil {
		return fmt.Errorf("failed to remove Public Ip Association for Nat Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed to wait for removal of Public Ip Association for Nat Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
