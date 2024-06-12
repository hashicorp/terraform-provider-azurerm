// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNATGatewayPublicIpAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNATGatewayPublicIpAssociationCreate,
		Read:   resourceNATGatewayPublicIpAssociationRead,
		Delete: resourceNATGatewayPublicIpAssociationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
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
				ValidateFunc: natgateways.ValidateNatGatewayID,
			},

			"public_ip_address_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidatePublicIPAddressID,
			},
		},
	}
}

func resourceNATGatewayPublicIpAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	publicIpAddressId, err := commonids.ParsePublicIPAddressID(d.Get("public_ip_address_id").(string))
	if err != nil {
		return err
	}

	natGatewayId, err := natgateways.ParseNatGatewayID(d.Get("nat_gateway_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(natGatewayId.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(natGatewayId.NatGatewayName, natGatewayResourceName)

	natGateway, err := client.Get(ctx, *natGatewayId, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(natGateway.HttpResponse) {
			return fmt.Errorf("%s was not found", natGatewayId)
		}
		return fmt.Errorf("retrieving %s: %+v", natGatewayId, err)
	}

	if natGateway.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", natGatewayId)
	}
	if natGateway.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", natGatewayId)
	}

	id := commonids.NewCompositeResourceID(natGatewayId, publicIpAddressId)

	publicIpAddresses := make([]natgateways.SubResource, 0)
	if natGateway.Model.Properties.PublicIPAddresses != nil {
		for _, existingPublicIPAddress := range *natGateway.Model.Properties.PublicIPAddresses {
			if existingPublicIPAddress.Id == nil {
				continue
			}

			if strings.EqualFold(*existingPublicIPAddress.Id, publicIpAddressId.ID()) {
				return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_association", id.ID())
			}

			publicIpAddresses = append(publicIpAddresses, existingPublicIPAddress)
		}
	}

	publicIpAddresses = append(publicIpAddresses, natgateways.SubResource{
		Id: pointer.To(publicIpAddressId.ID()),
	})
	natGateway.Model.Properties.PublicIPAddresses = &publicIpAddresses

	if err := client.CreateOrUpdateThenPoll(ctx, *natGatewayId, *natGateway.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", natGatewayId, err)
	}

	d.SetId(id.ID())

	return resourceNATGatewayPublicIpAssociationRead(d, meta)
}

func resourceNATGatewayPublicIpAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
	if err != nil {
		return err
	}

	natGateway, err := client.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(natGateway.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id.First)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if model := natGateway.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.PublicIPAddresses == nil {
				log.Printf("[DEBUG] %s doesn't have any Public IP's - removing from state!", id.First)
				d.SetId("")
				return nil
			}

			publicIPAddressId := ""
			for _, pip := range *props.PublicIPAddresses {
				if pip.Id == nil {
					continue
				}

				if strings.EqualFold(*pip.Id, id.Second.ID()) {
					publicIPAddressId = *pip.Id
					break
				}
			}

			if publicIPAddressId == "" {
				log.Printf("[DEBUG] Association between %s and %s was not found - removing from state", id.First, id.Second)
				d.SetId("")
				return nil
			}
		}
	}

	d.Set("nat_gateway_id", id.First.ID())
	d.Set("public_ip_address_id", id.Second.ID())

	return nil
}

func resourceNATGatewayPublicIpAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
	if err != nil {
		return err
	}

	locks.ByName(id.First.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(id.First.NatGatewayName, natGatewayResourceName)

	natGateway, err := client.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(natGateway.HttpResponse) {
			return fmt.Errorf("%s was not found", id.First)
		}
		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if natGateway.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id.First)
	}
	if natGateway.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
	}

	publicIpAddresses := make([]natgateways.SubResource, 0)
	if publicIPAddresses := natGateway.Model.Properties.PublicIPAddresses; publicIPAddresses != nil {
		for _, publicIPAddress := range *publicIPAddresses {
			if publicIPAddress.Id == nil {
				continue
			}

			if !strings.EqualFold(*publicIPAddress.Id, id.Second.ID()) {
				publicIpAddresses = append(publicIpAddresses, publicIPAddress)
			}
		}
	}
	natGateway.Model.Properties.PublicIPAddresses = &publicIpAddresses

	if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *natGateway.Model); err != nil {
		return fmt.Errorf("removing association between %s and %s: %+v", id.First, id.Second, err)
	}

	return nil
}
