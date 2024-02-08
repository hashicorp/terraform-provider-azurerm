// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayapi"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGatewayApi() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayApiCreate,
		Read:   resourceApiManagementGatewayApiRead,
		Delete: resourceApiManagementGatewayApiDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := gatewayapi.ParseGatewayApiID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiManagementGatewayApiV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"api_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: api.ValidateApiID,
			},
			"gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceApiManagementGatewayApiCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	apiName := getApiName(apiID.ApiId)

	id := gatewayapi.NewGatewayApiID(gatewayID.SubscriptionId, gatewayID.ResourceGroupName, gatewayID.ServiceName, gatewayID.GatewayId, apiName)
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
	params := gatewayapi.AssociationContract{}
	if _, err = client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceApiManagementGatewayApiRead(d, meta)
}

func resourceApiManagementGatewayApiRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayApisClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gatewayapi.ParseGatewayApiID(d.Id())
	if err != nil {
		return err
	}

	apiName := getApiName(id.ApiId)

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

func resourceApiManagementGatewayApiDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayApisClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gatewayapi.ParseGatewayApiID(d.Id())
	if err != nil {
		return err
	}

	name := getApiName(id.ApiId)

	newId := gatewayapi.NewGatewayApiID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GatewayId, name)
	if resp, err := client.Delete(ctx, newId); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("removing %s: %+v", newId, err)
		}
	}

	return nil
}
