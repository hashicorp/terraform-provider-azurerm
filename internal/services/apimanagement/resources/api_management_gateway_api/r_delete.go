// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_api

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayapi"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGatewayApiDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.GatewayApisClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gatewayapi.ParseGatewayApiID(d.Id())
	if err != nil {
		return err
	}

	name := helpers.GetApiName(id.ApiId)

	newId := gatewayapi.NewGatewayApiID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GatewayId, name)
	if resp, err := client.Delete(ctx, newId); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("removing %s: %+v", newId, err)
		}
	}

	return nil
}
