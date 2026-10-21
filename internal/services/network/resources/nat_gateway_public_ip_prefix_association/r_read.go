// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nat_gateway_public_ip_prefix_association

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNATGatewayPublicIpPrefixAssociationRead(d *pluginsdk.ResourceData, meta any) error {
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
