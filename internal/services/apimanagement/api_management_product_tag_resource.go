// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/producttag"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementProductTag() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementProductTagCreate,
		Read:   resourceApiManagementProductTagRead,
		Delete: resourceApiManagementProductTagDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := producttag.ParseProductTagID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
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
	client := meta.(*clients.Client).ApiManagement.ProductTagClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := producttag.NewProductTagID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_management_product_id").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.TagGetByProduct(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_product_tag", id.ID())
		}
	}

	resp, err := client.TagAssignToProduct(ctx, id)
	if err != nil {
		return fmt.Errorf(" creating product tag (id : %s): %+v", id, err)
	}
	if resp.Model != nil {
		d.SetId(pointer.From(resp.Model.Id))
	}

	return resourceApiManagementProductTagRead(d, meta)
}

func resourceApiManagementProductTagRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductTagClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := producttag.ParseProductTagID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.TagGetByProduct(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		productTagId, err := producttag.ParseProductTagID(*model.Id)
		if err != nil {
			return err
		}
		d.Set("api_management_product_id", productTagId.ProductId)
		d.Set("api_management_name", productTagId.ServiceName)
		d.Set("resource_group_name", productTagId.ResourceGroupName)
		d.Set("name", productTagId.TagId)
	}

	return nil
}

func resourceApiManagementProductTagDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductTagClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := producttag.ParseProductTagID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	resp, err := client.TagDetachFromProduct(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
