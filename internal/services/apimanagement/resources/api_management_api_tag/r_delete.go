// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiTagDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiTagClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apitag.ParseApiTagID(d.Id())
	if err != nil {
		return err
	}

	newId := apitag.NewApiTagID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.TagId)

	if _, err = client.TagDetachFromApi(ctx, newId); err != nil {
		return fmt.Errorf("detaching api tag %q: %+v", newId, err)
	}

	return nil
}
