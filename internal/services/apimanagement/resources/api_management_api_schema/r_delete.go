// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_schema

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apischema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiSchemaDelete(d *pluginsdk.ResourceData, meta any) error {
	client := apimanagementClient.From(meta).ApiSchemasClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := apischema.ParseApiSchemaID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, apischema.DeleteOperationOptions{Force: pointer.To(false)}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %s", *id, err)
		}
	}

	return nil
}
