// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_api

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayapi"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGatewayApiCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.GatewayApisClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiID, err := api.ParseApiID(d.Get("api_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_id`: %v", err)
	}

	gatewayID, err := gateway.ParseGatewayID(d.Get("gateway_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `gateway_id`: %v", err)
	}

	apiName := helpers.GetApiName(apiID.ApiId)

	id := gatewayapi.NewGatewayApiID(gatewayID.SubscriptionId, gatewayID.ResourceGroupName, gatewayID.ServiceName, gatewayID.GatewayId, apiName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		exists, err := client.GetEntityTag(ctx, id)
		if err != nil {
			if !response.WasStatusCode(exists.HttpResponse, http.StatusNoContent) {
				if !response.WasNotFound(exists.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", gatewayID, err)
				}
			}
		}

		if !response.WasNotFound(exists.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_gateway_api", id.ID())
		}
	}

	if _, err = client.CreateOrUpdate(ctx, id, gatewayapi.AssociationContract{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceApiManagementGatewayApiRead(d, meta)
}
