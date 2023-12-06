// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gateway"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayCreateUpdate,
		Read:   resourceApiManagementGatewayRead,
		Update: resourceApiManagementGatewayCreateUpdate,
		Delete: resourceApiManagementGatewayDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := gateway.ParseGatewayID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: apimanagementservice.ValidateServiceID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"location_data": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"city": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"district": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"region": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceApiManagementGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apimId, err := apimanagementservice.ParseServiceID(d.Get("api_management_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_id`: %v", err)
	}

	id := gateway.NewGatewayID(apimId.SubscriptionId, apimId.ResourceGroupName, apimId.ServiceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("making read request %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_gateway", id.ID())
		}
	}

	description := d.Get("description").(string)
	locationData := expandApiManagementGatewayLocationData(d.Get("location_data").([]interface{}))

	parameters := gateway.GatewayContract{
		Properties: &gateway.GatewayContractProperties{
			Description:  pointer.To(description),
			LocationData: locationData,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, gateway.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGatewayRead(d, meta)
}

func resourceApiManagementGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gateway.ParseGatewayID(d.Id())
	if err != nil {
		return err
	}
	apimId := apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request for %s: %+v", id, err)
	}

	d.Set("name", id.GatewayId)
	d.Set("api_management_id", apimId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("location_data", flattenApiManagementGatewayLocationData(props.LocationData))
		}
	}

	return nil
}

func resourceApiManagementGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gateway.ParseGatewayID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, gateway.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func expandApiManagementGatewayLocationData(input []interface{}) *gateway.ResourceLocationDataContract {
	if len(input) == 0 {
		return nil
	}

	locationData := gateway.ResourceLocationDataContract{}

	vs := input[0].(map[string]interface{})
	for k, v := range vs {
		switch k {
		case "name":
			locationData.Name = v.(string)
		case "city":
			locationData.City = pointer.To(v.(string))
		case "district":
			locationData.District = pointer.To(v.(string))
		case "region":
			locationData.CountryOrRegion = pointer.To(v.(string))
		}
	}

	return &locationData
}

func flattenApiManagementGatewayLocationData(input *gateway.ResourceLocationDataContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	locationData := map[string]interface{}{
		"name":     input.Name,
		"city":     pointer.From(input.City),
		"region":   pointer.From(input.CountryOrRegion),
		"district": pointer.From(input.District),
	}

	return []interface{}{locationData}
}
