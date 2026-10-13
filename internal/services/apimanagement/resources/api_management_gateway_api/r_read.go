// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayapi"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGatewayApiRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.GatewayApisClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gatewayapi.ParseGatewayApiID(d.Id())
	if err != nil {
		return err
	}

	apiName := helpers.GetApiName(id.ApiId)

	apiId := api.NewApiID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName)
	gatewayApiId := gatewayapi.NewGatewayApiID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GatewayId, apiName)
	resp, err := client.GetEntityTag(ctx, gatewayApiId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", gatewayApiId)
			d.SetId("")
			return nil
		}
		if response.WasStatusCode(resp.HttpResponse, http.StatusNoContent) {
			log.Printf("[DEBUG] %s returned with No Content status - bypassing and moving on!", gatewayApiId)
		} else {
			return fmt.Errorf("retrieving %s: %+v", gatewayApiId, err)
		}
	}
	if response.WasNotFound(resp.HttpResponse) {
		log.Printf("[DEBUG] %s was not found - removing from state!", *id)
		d.SetId("")
		return nil
	}
	gateway := gatewayapi.NewGatewayID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GatewayId)

	d.Set("api_id", apiId.ID())
	d.Set("gateway_id", gateway.ID())

	return nil
}
