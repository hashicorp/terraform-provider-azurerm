// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/publicipaddresses"
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

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceNATGatewayPublicIpAssociationCustomizeDiff),

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

func resourceNATGatewayPublicIpAssociationCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, meta any) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	rawNatGatewayId := d.GetRawConfig().AsValueMap()["nat_gateway_id"]
	if rawNatGatewayId.IsNull() || !rawNatGatewayId.IsKnown() {
		return nil
	}

	rawPublicIPAddressId := d.GetRawConfig().AsValueMap()["public_ip_address_id"]
	if rawPublicIPAddressId.IsNull() || !rawPublicIPAddressId.IsKnown() {
		return nil
	}

	natGatewayId, err := natgateways.ParseNatGatewayID(d.Get("nat_gateway_id").(string))
	if err != nil {
		return err
	}

	publicIpAddressId, err := commonids.ParsePublicIPAddressID(d.Get("public_ip_address_id").(string))
	if err != nil {
		return err
	}

	client := meta.(*clients.Client)
	natGateway, err := client.Network.NatGateways.Get(ctx, *natGatewayId, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(natGateway.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", natGatewayId, err)
	}
	if natGateway.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", natGatewayId)
	}

	publicIPAddress, err := client.Network.PublicIPAddresses.Get(ctx, *publicIpAddressId, publicipaddresses.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(publicIPAddress.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", publicIpAddressId, err)
	}
	if publicIPAddress.Model == nil || publicIPAddress.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model` or `properties` was nil", publicIpAddressId)
	}

	return validateNATGatewayPublicIpAssociation(natGateway.Model, publicIPAddress.Model)
}

func resourceNATGatewayPublicIpAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGateways
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

	publicIPAddress, err := meta.(*clients.Client).Network.PublicIPAddresses.Get(ctx, *publicIpAddressId, publicipaddresses.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(publicIPAddress.HttpResponse) {
			return fmt.Errorf("%s was not found", publicIpAddressId)
		}
		return fmt.Errorf("retrieving %s: %+v", publicIpAddressId, err)
	}
	if publicIPAddress.Model == nil || publicIPAddress.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model` or `properties` was nil", publicIpAddressId)
	}

	isIPv6 := natGatewayPublicIpAssociationIsIPv6(publicIPAddress.Model)
	if err := validateNATGatewayPublicIpAssociation(natGateway.Model, publicIPAddress.Model); err != nil {
		return err
	}

	id := commonids.NewCompositeResourceID(natGatewayId, publicIpAddressId)

	publicIpAddresses := pointer.From(natGateway.Model.Properties.PublicIPAddresses)
	if isIPv6 {
		publicIpAddresses = pointer.From(natGateway.Model.Properties.PublicIPAddressesV6)
	}
	for _, existingPublicIPAddress := range publicIpAddresses {
		if strings.EqualFold(pointer.From(existingPublicIPAddress.Id), publicIpAddressId.ID()) {
			return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_association", id.ID())
		}
	}

	publicIpAddresses = append(publicIpAddresses, natgateways.SubResource{
		Id: pointer.To(publicIpAddressId.ID()),
	})
	if isIPv6 {
		natGateway.Model.Properties.PublicIPAddressesV6 = pointer.To(publicIpAddresses)
	} else {
		natGateway.Model.Properties.PublicIPAddresses = pointer.To(publicIpAddresses)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *natGatewayId, *natGateway.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", natGatewayId, err)
	}

	d.SetId(id.ID())

	return resourceNATGatewayPublicIpAssociationRead(d, meta)
}

func resourceNATGatewayPublicIpAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	if model := natGateway.Model; model != nil && model.Properties != nil {
		if !natGatewayPublicIpAssociationExists(model.Properties, id.Second.ID()) {
			log.Printf("[DEBUG] Association between %s and %s was not found - removing from state", id.First, id.Second)
			d.SetId("")
			return nil
		}
	} else {
		return fmt.Errorf("retrieving %s: `model` or `properties` was nil", id.First)
	}

	d.Set("nat_gateway_id", id.First.ID())
	d.Set("public_ip_address_id", id.Second.ID())

	return nil
}

func resourceNATGatewayPublicIpAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &natgateways.NatGatewayId{}, &commonids.PublicIPAddressId{})
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

	if !removeNATGatewayPublicIpAssociation(natGateway.Model.Properties, id.Second.ID()) {
		return nil
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *natGateway.Model); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func validateNATGatewayPublicIpAssociation(natGateway *natgateways.NatGateway, publicIPAddress *publicipaddresses.PublicIPAddress) error {
	isIPv6 := natGatewayPublicIpAssociationIsIPv6(publicIPAddress)
	natGatewaySku := pointer.From(pointer.From(natGateway.Sku).Name)
	publicIPAddressSku := pointer.From(pointer.From(publicIPAddress.Sku).Name)

	if isIPv6 {
		if natGatewaySku != natgateways.NatGatewaySkuNameStandardVTwo || publicIPAddressSku != publicipaddresses.PublicIPAddressSkuNameStandardVTwo {
			return errors.New("`nat_gateway_id` must reference a NAT Gateway with SKU `StandardV2` and `public_ip_address_id` must reference an `IPv6` Public IP Address with SKU `StandardV2` when `public_ip_address_id` references an `IPv6` Public IP Address")
		}
	}

	if natGatewaySku == natgateways.NatGatewaySkuNameStandard && publicIPAddressSku == publicipaddresses.PublicIPAddressSkuNameStandardVTwo {
		return errors.New("`public_ip_address_id` must reference a Public IP Address with SKU `Standard` when `nat_gateway_id` references a NAT Gateway with SKU `Standard`")
	}

	if natGatewaySku == natgateways.NatGatewaySkuNameStandardVTwo && publicIPAddressSku != publicipaddresses.PublicIPAddressSkuNameStandardVTwo {
		return errors.New("`public_ip_address_id` must reference a Public IP Address with SKU `StandardV2` when `nat_gateway_id` references a NAT Gateway with SKU `StandardV2`")
	}

	return nil
}

func natGatewayPublicIpAssociationIsIPv6(publicIPAddress *publicipaddresses.PublicIPAddress) bool {
	return pointer.From(publicIPAddress.Properties.PublicIPAddressVersion) == publicipaddresses.IPVersionIPvSix
}

func natGatewayPublicIpAssociationExists(properties *natgateways.NatGatewayPropertiesFormat, publicIPAddressId string) bool {
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
