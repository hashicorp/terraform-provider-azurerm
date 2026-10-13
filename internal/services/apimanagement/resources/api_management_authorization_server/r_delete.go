// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_authorization_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementAuthorizationServerDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.AuthorizationServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := authorizationserver.ParseAuthorizationServerID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, authorizationserver.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %s", *id, err)
		}
	}

	return nil
}
