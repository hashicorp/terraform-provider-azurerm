// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation_tag

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperationtag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiOperationTagDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationTagClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apioperationtag.ParseOperationTagID(d.Id())
	if err != nil {
		return err
	}

	apiName := helpers.GetApiName(id.ApiId)

	newId := apioperationtag.NewOperationTagID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, apiName, id.OperationId, id.TagId)
	if _, err = client.TagDetachFromOperation(ctx, newId); err != nil {
		return fmt.Errorf("deleting %q: %+v", newId, err)
	}

	return nil
}
