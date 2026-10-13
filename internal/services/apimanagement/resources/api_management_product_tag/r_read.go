// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_product_tag

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/producttag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementProductTagRead(d *pluginsdk.ResourceData, meta any) error {
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
