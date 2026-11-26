// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation_tag

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationtag"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiOperationTagUpdate(d *pluginsdk.ResourceData, meta any) error {
	tagClient := meta.(*clients.Client).ApiManagement.TagClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apioperationtag.ParseOperationTagID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("display_name") {
		tagId := tag.NewTagID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.TagId)
		if err := setApiManagementTagDisplayName(ctx, tagClient, tagId, d.Get("display_name").(string)); err != nil {
			return err
		}
	}

	return resourceApiManagementApiOperationTagRead(d, meta)
}
