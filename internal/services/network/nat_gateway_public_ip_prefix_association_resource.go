// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/publicipprefixes"
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
	client := meta.(*clients.Client).Network.NatGateways
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

	locks.ByID(natGatewayId.ID())
	defer locks.UnlockByID(natGatewayId.ID())

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

	publicIpPrefix, err := meta.(*clients.Client).Network.PublicIPPrefixes.Get(ctx, *publicIpPrefixId, publicipprefixes.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(publicIpPrefix.HttpResponse) {
			return fmt.Errorf("%s was not found", publicIpPrefixId)
		}
		return fmt.Errorf("retrieving %s: %+v", publicIpPrefixId, err)
	}
	if publicIpPrefix.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", publicIpPrefixId)
	}
	if publicIpPrefix.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", publicIpPrefixId)
	}

	isIPv6 := pointer.From(publicIpPrefix.Model.Properties.PublicIPAddressVersion) == publicipprefixes.IPVersionIPvSix
	id := commonids.NewCompositeResourceID(natGatewayId, publicIpPrefixId)

	gatewayProperties := natGateway.Model.Properties
	publicIpPrefixes := pointer.From(gatewayProperties.PublicIPPrefixes)
	if isIPv6 {
		publicIpPrefixes = pointer.From(gatewayProperties.PublicIPPrefixesV6)
	}
	for _, existingPublicIPPrefix := range publicIpPrefixes {
		if strings.EqualFold(pointer.From(existingPublicIPPrefix.Id), publicIpPrefixId.ID()) {
			return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_prefix_association", id.ID())
		}
	}

	publicIpPrefixes = append(publicIpPrefixes, natgateways.SubResource{
		Id: pointer.To(publicIpPrefixId.ID()),
	})
	if isIPv6 {
		gatewayProperties.PublicIPPrefixesV6 = pointer.To(publicIpPrefixes)
	} else {
		gatewayProperties.PublicIPPrefixes = pointer.To(publicIpPrefixes)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *natGatewayId, *natGateway.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", natGatewayId, err)
	}

	d.SetId(id.ID())

	return resourceNATGatewayPublicIpPrefixAssociationRead(d, meta)
}

func resourceNATGatewayPublicIpPrefixAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGateways
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

	if natGateway.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id.First)
	}
	if natGateway.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
	}
	if !natGatewayPublicIpPrefixAssociationExists(natGateway.Model.Properties, id.Second.ID()) {
		log.Printf("[DEBUG] Association between %s and %s was not found - removing from state", id.First, id.Second)
		d.SetId("")
		return nil
	}

	d.Set("nat_gateway_id", id.First.ID())
	d.Set("public_ip_prefix_id", id.Second.ID())

	return nil
}

func resourceNATGatewayPublicIpPrefixAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &natgateways.NatGatewayId{}, &publicipprefixes.PublicIPPrefixId{})
	if err != nil {
		return err
	}

	locks.ByID(id.First.ID())
	defer locks.UnlockByID(id.First.ID())

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

	if !removeNATGatewayPublicIpPrefixAssociation(natGateway.Model.Properties, id.Second.ID()) {
		return nil
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *natGateway.Model); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func natGatewayPublicIpPrefixAssociationExists(properties *natgateways.NatGatewayPropertiesFormat, publicIpPrefixId string) bool {
	if properties == nil {
		return false
	}

	for _, publicIpPrefix := range pointer.From(properties.PublicIPPrefixes) {
		if strings.EqualFold(pointer.From(publicIpPrefix.Id), publicIpPrefixId) {
			return true
		}
	}

	for _, publicIpPrefix := range pointer.From(properties.PublicIPPrefixesV6) {
		if strings.EqualFold(pointer.From(publicIpPrefix.Id), publicIpPrefixId) {
			return true
		}
	}

	return false
}

func removeNATGatewayPublicIpPrefixAssociation(properties *natgateways.NatGatewayPropertiesFormat, publicIpPrefixId string) bool {
	if properties == nil {
		return false
	}

	removed := false
	updatedIPv4Prefixes := make([]natgateways.SubResource, 0)
	for _, publicIpPrefix := range pointer.From(properties.PublicIPPrefixes) {
		if strings.EqualFold(pointer.From(publicIpPrefix.Id), publicIpPrefixId) {
			removed = true
			continue
		}

		updatedIPv4Prefixes = append(updatedIPv4Prefixes, publicIpPrefix)
	}
	properties.PublicIPPrefixes = pointer.To(updatedIPv4Prefixes)

	updatedIPv6Prefixes := make([]natgateways.SubResource, 0)
	for _, publicIpPrefix := range pointer.From(properties.PublicIPPrefixesV6) {
		if strings.EqualFold(pointer.From(publicIpPrefix.Id), publicIpPrefixId) {
			removed = true
			continue
		}

		updatedIPv6Prefixes = append(updatedIPv6Prefixes, publicIpPrefix)
	}
	properties.PublicIPPrefixesV6 = pointer.To(updatedIPv6Prefixes)

	return removed
}
