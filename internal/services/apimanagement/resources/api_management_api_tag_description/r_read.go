// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag_description

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitagdescription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiTagDescriptionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apitagdescription.ParseTagDescriptionID(d.Id())
	if err != nil {
		return err
	}

	apiName := helpers.GetApiName(id.ApiId)

	newId := apitagdescription.NewTagDescriptionID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName, id.TagDescriptionId)
	resp, err := client.Get(ctx, newId)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", newId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", newId, err)
	}

	apiTagId := apitag.NewApiTagID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName, id.TagDescriptionId)

	d.Set("api_tag_id", apiTagId.ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("external_documentation_url", pointer.From(props.ExternalDocsURL))
			d.Set("external_documentation_description", pointer.From(props.ExternalDocsDescription))
		}
	}

	return nil
}
