// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementProductTag() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementProductTagCreate,
		Read:   resourceApiManagementProductTagRead,
		Delete: resourceApiManagementProductTagDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ProductTagID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_management_product_id": schemaz.SchemaApiManagementChildName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"name": schemaz.SchemaApiManagementChildName(),
		},
	}
}

func resourceApiManagementProductTagCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.TagClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewProductTagID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_management_product_id").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.GetByProduct(ctx, id.ResourceGroup, id.ServiceName, id.ProductName, id.TagName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_product_tag", id.ID())
		}
	}

	resp, err := client.AssignToProduct(ctx, id.ResourceGroup, id.ServiceName, id.ProductName, id.TagName)
	if err != nil {
		return fmt.Errorf(" creating product tag (id : %s): %+v", id, err)
	}
	d.SetId(*resp.ID)

	return resourceApiManagementProductTagRead(d, meta)
}

func resourceApiManagementProductTagRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductTagID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetByProduct(ctx, id.ResourceGroup, id.ServiceName, id.ProductName, id.TagName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	productTagId, err := parse.ProductTagID(*resp.ID)
	if err != nil {
		return err
	}

	d.Set("api_management_product_id", productTagId.ProductName)
	d.Set("api_management_name", productTagId.ServiceName)
	d.Set("resource_group_name", productTagId.ResourceGroup)
	d.Set("name", productTagId.TagName)

	return nil
}

func resourceApiManagementProductTagDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductTagID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	resp, err := client.DetachFromProduct(ctx, id.ResourceGroup, id.ServiceName, id.ProductName, id.TagName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
