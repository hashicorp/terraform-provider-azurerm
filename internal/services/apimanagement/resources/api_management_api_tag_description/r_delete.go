// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag_description

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitagdescription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiTagDescriptionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagDescriptionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apitagdescription.ParseTagDescriptionID(d.Id())
	if err != nil {
		return err
	}

	name := helpers.GetApiName(id.ApiId)

	newId := apitagdescription.NewTagDescriptionID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, name, id.TagDescriptionId)
	resp, err := client.Delete(ctx, newId, apitagdescription.DeleteOperationOptions{})
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", newId, err)
		}
	}

	return nil
}
