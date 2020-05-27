package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNATGatewayPublicIpAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNATGatewayPublicIpAssociationCreate,
		Read:   resourceArmNATGatewayPublicIpAssociationRead,
		Delete: resourceArmNATGatewayPublicIpAssociationDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.NatGatewayPublicIPAddressAssociationID(id)
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
				ValidateFunc: validate.PublicIPAddressID,
			},
		},
	}
}

func resourceArmNATGatewayPublicIpAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for NAT Gateway <-> Public IP Association creation.")
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
			return fmt.Errorf("NAT Gateway %q (Resource Group %q) was not found.", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup)
		}
		return fmt.Errorf("failed to retrieve NAT Gateway %q (Resource Group %q): %+v", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup, err)
	}

	id := fmt.Sprintf("%s|%s", *natGateway.ID, publicIpAddressId)
	publicIpAddresses := make([]network.SubResource, 0)
	if natGateway.PublicIPAddresses != nil {
		for _, existingPublicIPAddress := range *natGateway.PublicIPAddresses {
			if existingPublicIPAddress.ID == nil {
				continue
			}

			if strings.EqualFold(*existingPublicIPAddress.ID, publicIpAddressId) {
				return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_association", id)
			}

			publicIpAddresses = append(publicIpAddresses, existingPublicIPAddress)
		}
	}

	publicIpAddresses = append(publicIpAddresses, network.SubResource{
		ID: utils.String(publicIpAddressId),
	})
	natGateway.PublicIPAddresses = &publicIpAddresses

	future, err := client.CreateOrUpdate(ctx, parsedNatGatewayId.ResourceGroup, parsedNatGatewayId.Name, natGateway)
	if err != nil {
		return fmt.Errorf("failed to update Public IP Association for NAT Gateway %q (Resource Group %q): %+v", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed to wait for completion of Public IP Association for NAT Gateway %q (Resource Group %q): %+v", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup, err)
	}

	d.SetId(id)

	return resourceArmNATGatewayPublicIpAssociationRead(d, meta)
}

func resourceArmNATGatewayPublicIpAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayPublicIPAddressAssociationID(d.Id())
	if err != nil {
		return err
	}

	natGateway, err := client.Get(ctx, id.NatGateway.ResourceGroup, id.NatGateway.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(natGateway.Response) {
			log.Printf("[DEBUG] NAT Gateway %q (Resource Group %q) could not be found - removing from state!", id.NatGateway.Name, id.NatGateway.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to retrieve NAT Gateway %q (Resource Group %q): %+v", id.NatGateway.Name, id.NatGateway.ResourceGroup, err)
	}

	if natGateway.NatGatewayPropertiesFormat == nil {
		return fmt.Errorf("`properties` was nil for NAT Gateway %q (Resource Group %q)", id.NatGateway.Name, id.NatGateway.ResourceGroup)
	}
	props := *natGateway.NatGatewayPropertiesFormat

	if props.PublicIPAddresses == nil {
		log.Printf("[DEBUG] NAT Gateway %q (Resource Group %q) doesn't have any Public IP's - removing from state!", id.NatGateway.Name, id.NatGateway.ResourceGroup)
		d.SetId("")
		return nil
	}

	publicIPAddressId := ""
	for _, pip := range *props.PublicIPAddresses {
		if pip.ID == nil {
			continue
		}

		if strings.EqualFold(*pip.ID, id.PublicIPAddressID) {
			publicIPAddressId = *pip.ID
			break
		}
	}

	if publicIPAddressId == "" {
		log.Printf("[DEBUG] Association between NAT Gateway %q (Resource Group %q) and Public IP Address %q was not found - removing from state", id.NatGateway.Name, id.NatGateway.ResourceGroup, id.PublicIPAddressID)
		d.SetId("")
		return nil
	}

	d.Set("nat_gateway_id", natGateway.ID)
	d.Set("public_ip_address_id", publicIPAddressId)

	return nil
}

func resourceArmNATGatewayPublicIpAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayPublicIPAddressAssociationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NatGateway.Name, natGatewayResourceName)
	defer locks.UnlockByName(id.NatGateway.Name, natGatewayResourceName)

	natGateway, err := client.Get(ctx, id.NatGateway.ResourceGroup, id.NatGateway.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(natGateway.Response) {
			return fmt.Errorf("NAT Gateway %q (Resource Group %q) was not found", id.NatGateway.Name, id.NatGateway.ResourceGroup)
		}

		return fmt.Errorf("retrieving NAT Gateway %q (Resource Group %q): %+v", id.NatGateway.Name, id.NatGateway.ResourceGroup, err)
	}
	if natGateway.NatGatewayPropertiesFormat == nil {
		return fmt.Errorf("retrieving NAT Gateway %q (Resource Group %q): `properties` was nil", id.NatGateway.Name, id.NatGateway.ResourceGroup)
	}

	publicIpAddresses := make([]network.SubResource, 0)
	if publicIPAddresses := natGateway.NatGatewayPropertiesFormat.PublicIPAddresses; publicIPAddresses != nil {
		for _, publicIPAddress := range *publicIPAddresses {
			if publicIPAddress.ID == nil {
				continue
			}

			if !strings.EqualFold(*publicIPAddress.ID, id.PublicIPAddressID) {
				publicIpAddresses = append(publicIpAddresses, publicIPAddress)
			}
		}
	}
	natGateway.NatGatewayPropertiesFormat.PublicIPAddresses = &publicIpAddresses

	future, err := client.CreateOrUpdate(ctx, id.NatGateway.ResourceGroup, id.NatGateway.Name, natGateway)
	if err != nil {
		return fmt.Errorf("removing association between NAT Gateway %q (Resource Group %q) and Public IP Address %q: %+v", id.NatGateway.Name, id.NatGateway.ResourceGroup, id.PublicIPAddressID, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for association between Public IP ID %q for NAT Gateway %q (Resource Group %q) to be removed: %+v", id.PublicIPAddressID, id.NatGateway.Name, id.NatGateway.ResourceGroup, err)
	}

	return nil
}
