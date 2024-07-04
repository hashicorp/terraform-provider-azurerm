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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/publicipprefixes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNATGatewayPublicIpPrefixAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNATGatewayPublicIpPrefixAssociationCreate,
		Read:   resourceNATGatewayPublicIpPrefixAssociationRead,
		Delete: resourceNATGatewayPublicIpPrefixAssociationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
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

			"public_ip_prefix_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
			},
		},
	}
}

func resourceNATGatewayPublicIpPrefixAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	publicIpPrefixId, err := publicipprefixes.ParsePublicIPPrefixID(d.Get("public_ip_prefix_id").(string))
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

	id := commonids.NewCompositeResourceID(natGatewayId, publicIpPrefixId)

	publicIpPrefixes := make([]natgateways.SubResource, 0)
	if natGateway.Model.Properties.PublicIPPrefixes != nil {
		for _, existingPublicIPPrefix := range *natGateway.Model.Properties.PublicIPPrefixes {
			if existingPublicIPPrefix.Id == nil {
				continue
			}

			if strings.EqualFold(*existingPublicIPPrefix.Id, publicIpPrefixId.ID()) {
				return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_prefix_association", id.ID())
			}

			publicIpPrefixes = append(publicIpPrefixes, existingPublicIPPrefix)
		}
	}

	publicIpPrefixes = append(publicIpPrefixes, natgateways.SubResource{
		Id: pointer.To(publicIpPrefixId.ID()),
	})
	natGateway.Model.Properties.PublicIPPrefixes = &publicIpPrefixes

	if err := client.CreateOrUpdateThenPoll(ctx, *natGatewayId, *natGateway.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", natGatewayId, err)
	}

	d.SetId(id.ID())

	return resourceNATGatewayPublicIpPrefixAssociationRead(d, meta)
}

func resourceNATGatewayPublicIpPrefixAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
	if err != nil {
		return err
	}

	natGateway, err := client.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(natGateway.HttpResponse) {
			log.Printf("[DEBUG] %s could not be found - removing from state!", id.First)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if model := natGateway.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.PublicIPPrefixes == nil {
				log.Printf("[DEBUG] %s doesn't have any Public IP Prefixes - removing from state!", id.First)
				d.SetId("")
				return nil
			}

			publicIPPrefixId := ""
			for _, pipp := range *props.PublicIPPrefixes {
				if pipp.Id == nil {
					continue
				}

				if strings.EqualFold(*pipp.Id, id.Second.ID()) {
					publicIPPrefixId = *pipp.Id
					break
				}
			}

			if publicIPPrefixId == "" {
				log.Printf("[DEBUG] Association between %s and %s was not found - removing from state", id.First, id.Second)
				d.SetId("")
				return nil
			}
		}
	}

	d.Set("nat_gateway_id", id.First.ID())
	d.Set("public_ip_prefix_id", id.Second.ID())

	return nil
}

func resourceNATGatewayPublicIpPrefixAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
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

	publicIpPrefixes := make([]natgateways.SubResource, 0)
	if publicIPPrefixes := natGateway.Model.Properties.PublicIPPrefixes; publicIPPrefixes != nil {
		for _, publicIPPrefix := range *publicIPPrefixes {
			if publicIPPrefix.Id == nil {
				continue
			}

			if !strings.EqualFold(*publicIPPrefix.Id, id.Second.ID()) {
				publicIpPrefixes = append(publicIpPrefixes, publicIPPrefix)
			}
		}
	}
	natGateway.Model.Properties.PublicIPPrefixes = &publicIpPrefixes

	if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *natGateway.Model); err != nil {
		return fmt.Errorf("removing association between %s and %s: %+v", id.First, id.Second, err)
	}

	return nil
}
