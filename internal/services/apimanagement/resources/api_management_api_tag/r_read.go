// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiTagRead(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.ApiTagClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apitag.ParseApiTagID(d.Id())
	if err != nil {
		return err
	}

	apiId := api.NewApiID(subscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId)
	tagId := apitag.NewApiTagID(subscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.TagId)

	resp, err := client.TagGetByApi(ctx, tagId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", tagId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", tagId, err)
	}

	d.Set("api_id", apiId.ID())
	d.Set("name", id.TagId)

	return nil
}
