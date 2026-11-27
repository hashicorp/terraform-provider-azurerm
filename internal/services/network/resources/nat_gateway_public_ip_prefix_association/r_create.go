// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package nat_gateway_public_ip_prefix_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNATGatewayPublicIpPrefixAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).NatGateways
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
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

	publicIpPrefix, err := networkClient.From(meta).PublicIPPrefixes.Get(ctx, *publicIpPrefixId, publicipprefixes.DefaultGetOperationOptions())
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
	exists := false
	for _, existingPublicIPPrefix := range publicIpPrefixes {
		if strings.EqualFold(pointer.From(existingPublicIPPrefix.Id), publicIpPrefixId.ID()) {
			exists = true
			if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_prefix_association", id.ID())
			}
		}
	}

	if !exists {
		publicIpPrefixes = append(publicIpPrefixes, natgateways.SubResource{
			Id: pointer.To(publicIpPrefixId.ID()),
		})
	}
	if isIPv6 {
		gatewayProperties.PublicIPPrefixesV6 = pointer.To(publicIpPrefixes)
	} else {
		gatewayProperties.PublicIPPrefixes = pointer.To(publicIpPrefixes)
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, *natGatewayId, *natGateway.Model, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("updating %s: %+v", natGatewayId, err)
	}

	d.SetId(id.ID())

	return resourceNATGatewayPublicIpPrefixAssociationRead(d, meta)
}
