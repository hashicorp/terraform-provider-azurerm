// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_account

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbAccountDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseDatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DatabaseAccountsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	// the SDK now will return a `WasNotFound` response even when still deleting
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"NotFound"},
		MinTimeout: 30 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
		Refresh: func() (any, string, error) {
			resp, err2 := client.DatabaseAccountsGet(ctx, *id)
			if err2 != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return resp, "NotFound", nil
				}
				return nil, "", err2
			}

			return resp, "Deleting", nil
		},
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}
