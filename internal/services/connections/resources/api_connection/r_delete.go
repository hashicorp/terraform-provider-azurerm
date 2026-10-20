// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiConnectionDelete(d *schema.ResourceData, meta any) error {
	client := meta.(*clients.Client).Connections.ConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
