// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_global_schema

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGlobalSchemaDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.GlobalSchemaClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schema.ParseSchemaID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.GlobalSchemaDelete(ctx, *id, schema.DefaultGlobalSchemaDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %s", *id, err)
	}

	return nil
}
