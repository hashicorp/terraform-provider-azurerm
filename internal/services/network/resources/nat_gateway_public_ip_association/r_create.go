// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nat_gateway_public_ip_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNATGatewayPublicIpAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
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
	if publicIPAddress.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", publicIpAddressId)
	}
	if publicIPAddress.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", publicIpAddressId)
	}

	isIPv6 := pointer.From(publicIPAddress.Model.Properties.PublicIPAddressVersion) == publicipaddresses.IPVersionIPvSix
	id := commonids.NewCompositeResourceID(natGatewayId, publicIpAddressId)

	gatewayProps := natGateway.Model.Properties
	publicIpAddresses := pointer.From(gatewayProps.PublicIPAddresses)
	if isIPv6 {
		publicIpAddresses = pointer.From(gatewayProps.PublicIPAddressesV6)
	}

	exists := false
	for _, existingPublicIPAddress := range publicIpAddresses {
		if strings.EqualFold(pointer.From(existingPublicIPAddress.Id), publicIpAddressId.ID()) {
			exists = true
			if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				return tf.ImportAsExistsError("azurerm_nat_gateway_public_ip_association", id.ID())
			}
		}
	}

	if !exists {
		publicIpAddresses = append(publicIpAddresses, natgateways.SubResource{
			Id: pointer.To(publicIpAddressId.ID()),
		})
	}
	if isIPv6 {
		gatewayProps.PublicIPAddressesV6 = pointer.To(publicIpAddresses)
	} else {
		gatewayProps.PublicIPAddresses = pointer.To(publicIpAddresses)
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, *natGatewayId, *natGateway.Model, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("updating %s: %+v", natGatewayId, err)
	}

	d.SetId(id.ID())

	return resourceNATGatewayPublicIpAssociationRead(d, meta)
}
