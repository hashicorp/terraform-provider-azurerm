// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package nat_gateway

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNatGatewayRead(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).NatGateways
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := natgateways.ParseNatGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return resourceNatGatewayFlatten(d, id, resp.Model)
}
