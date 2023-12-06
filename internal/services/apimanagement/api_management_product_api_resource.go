// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/productapi"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementProductApi() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementProductApiCreate,
		Read:   resourceApiManagementProductApiRead,
		Delete: resourceApiManagementProductApiDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := productapi.ParseProductApiID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_name": schemaz.SchemaApiManagementApiName(),

			"product_id": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),
		},
	}
}

func resourceApiManagementProductApiCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductApisClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := productapi.NewProductApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("product_id").(string), d.Get("api_name").(string))

	exists, err := client.CheckEntityExists(ctx, id)
	if err != nil {
		if !response.WasNotFound(exists.HttpResponse) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(exists.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_api_management_product_api", id.ID())
	}

	if _, err := client.CreateOrUpdate(ctx, id); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementProductApiRead(d, meta)
}

func resourceApiManagementProductApiRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductApisClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := productapi.ParseProductApiID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.CheckEntityExists(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("api_name", getApiName(id.ApiId))
	d.Set("product_id", id.ProductId)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	return nil
}

func resourceApiManagementProductApiDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductApisClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := productapi.ParseProductApiID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
