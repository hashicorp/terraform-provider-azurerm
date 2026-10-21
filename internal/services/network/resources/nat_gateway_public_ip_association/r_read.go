// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nat_gateway_public_ip_association

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNATGatewayPublicIpAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NatGateways
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

	if natGateway.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id.First)
	}
	if natGateway.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
	}
	if !natGatewayPublicIpAssociationExists(natGateway.Model.Properties, id.Second.ID()) {
		log.Printf("[DEBUG] Association between %s and %s was not found - removing from state", id.First, id.Second)
		d.SetId("")
		return nil
	}

	d.Set("nat_gateway_id", id.First.ID())
	d.Set("public_ip_address_id", id.Second.ID())

	return nil
}

func natGatewayPublicIpAssociationExists(properties *natgateways.NatGatewayPropertiesFormat, publicIPAddressId string) bool {
	if properties == nil {
		return false
	}

	for _, publicIPAddress := range pointer.From(properties.PublicIPAddresses) {
		if strings.EqualFold(pointer.From(publicIPAddress.Id), publicIPAddressId) {
			return true
		}
	}

	for _, publicIPAddress := range pointer.From(properties.PublicIPAddressesV6) {
		if strings.EqualFold(pointer.From(publicIPAddress.Id), publicIPAddressId) {
			return true
		}
	}

	return false
}

func removeNATGatewayPublicIpAssociation(properties *natgateways.NatGatewayPropertiesFormat, publicIPAddressId string) bool {
	if properties == nil {
		return false
	}

	removed := false
	updatedIPv4Addresses := make([]natgateways.SubResource, 0)
	for _, publicIPAddress := range pointer.From(properties.PublicIPAddresses) {
		if strings.EqualFold(pointer.From(publicIPAddress.Id), publicIPAddressId) {
			removed = true
			continue
		}

		updatedIPv4Addresses = append(updatedIPv4Addresses, publicIPAddress)
	}
	properties.PublicIPAddresses = pointer.To(updatedIPv4Addresses)

	updatedIPv6Addresses := make([]natgateways.SubResource, 0)
	for _, publicIPAddress := range pointer.From(properties.PublicIPAddressesV6) {
		if strings.EqualFold(pointer.From(publicIPAddress.Id), publicIPAddressId) {
			removed = true
			continue
		}

		updatedIPv6Addresses = append(updatedIPv6Addresses, publicIPAddress)
	}
	properties.PublicIPAddressesV6 = pointer.To(updatedIPv6Addresses)

	return removed
}
