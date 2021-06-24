package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceNATGatewayPublicIpPrefixAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNATGatewayPublicIpPrefixAssociationCreate,
		Read:   resourceNATGatewayPublicIpPrefixAssociationRead,
		Delete: resourceNATGatewayPublicIpPrefixAssociationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NatGatewayPublicIPPrefixAssociationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"nat_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NatGatewayID,
			},

			"public_ip_prefix_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PublicIpPrefixID,
			},
		},
	}
}

func resourceNATGatewayPublicIpPrefixAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for NAT Gateway <-> Public IP Prefix Association creation.")
	natGatewayId := d.Get("nat_gateway_id").(string)
	publicIpPrefixId := d.Get("public_ip_prefix_id").(string)
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

	id := fmt.Sprintf("%s|%s", *natGateway.ID, publicIpPrefixId)
	publicIpPrefixes := make([]network.SubResource, 0)
	if natGateway.PublicIPPrefixes != nil {
		for _, existingPublicIPPrefix := range *natGateway.PublicIPPrefixes {
			if existingPublicIPPrefix.ID == nil {
				continue
			}

			if strings.EqualFold(*existingPublicIPPrefix.ID, publicIpPrefixId) {
				return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_prefix_association", id)
			}

			publicIpPrefixes = append(publicIpPrefixes, existingPublicIPPrefix)
		}
	}

	publicIpPrefixes = append(publicIpPrefixes, network.SubResource{
		ID: utils.String(publicIpPrefixId),
	})
	natGateway.PublicIPPrefixes = &publicIpPrefixes

	future, err := client.CreateOrUpdate(ctx, parsedNatGatewayId.ResourceGroup, parsedNatGatewayId.Name, natGateway)
	if err != nil {
		return fmt.Errorf("failed to update Public IP Prefix Association for NAT Gateway %q (Resource Group %q): %+v", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed to wait for completion of Public IP Prefix Association for NAT Gateway %q (Resource Group %q): %+v", parsedNatGatewayId.Name, parsedNatGatewayId.ResourceGroup, err)
	}

	d.SetId(id)

	return resourceNATGatewayPublicIpPrefixAssociationRead(d, meta)
}

func resourceNATGatewayPublicIpPrefixAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayPublicIPPrefixAssociationID(d.Id())
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

	if props.PublicIPPrefixes == nil {
		log.Printf("[DEBUG] NAT Gateway %q (Resource Group %q) doesn't have any Public IP Prefixes - removing from state!", id.NatGateway.Name, id.NatGateway.ResourceGroup)
		d.SetId("")
		return nil
	}

	publicIPPrefixId := ""
	for _, pipp := range *props.PublicIPPrefixes {
		if pipp.ID == nil {
			continue
		}

		if strings.EqualFold(*pipp.ID, id.PublicIPPrefixID) {
			publicIPPrefixId = *pipp.ID
			break
		}
	}

	if publicIPPrefixId == "" {
		log.Printf("[DEBUG] Association between NAT Gateway %q (Resource Group %q) and Public IP Prefix %q was not found - removing from state", id.NatGateway.Name, id.NatGateway.ResourceGroup, id.PublicIPPrefixID)
		d.SetId("")
		return nil
	}

	d.Set("nat_gateway_id", natGateway.ID)
	d.Set("public_ip_prefix_id", publicIPPrefixId)

	return nil
}

func resourceNATGatewayPublicIpPrefixAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayPublicIPPrefixAssociationID(d.Id())
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

	publicIpPrefixes := make([]network.SubResource, 0)
	if publicIPPrefixes := natGateway.NatGatewayPropertiesFormat.PublicIPPrefixes; publicIPPrefixes != nil {
		for _, publicIPPrefix := range *publicIPPrefixes {
			if publicIPPrefix.ID == nil {
				continue
			}

			if !strings.EqualFold(*publicIPPrefix.ID, id.PublicIPPrefixID) {
				publicIpPrefixes = append(publicIpPrefixes, publicIPPrefix)
			}
		}
	}
	natGateway.NatGatewayPropertiesFormat.PublicIPPrefixes = &publicIpPrefixes

	future, err := client.CreateOrUpdate(ctx, id.NatGateway.ResourceGroup, id.NatGateway.Name, natGateway)
	if err != nil {
		return fmt.Errorf("removing association between NAT Gateway %q (Resource Group %q) and Public IP Prefix %q: %+v", id.NatGateway.Name, id.NatGateway.ResourceGroup, id.PublicIPPrefixID, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for association between Public IP Prefix ID %q for NAT Gateway %q (Resource Group %q) to be removed: %+v", id.PublicIPPrefixID, id.NatGateway.Name, id.NatGateway.ResourceGroup, err)
	}

	return nil
}
