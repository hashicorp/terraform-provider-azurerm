// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_key

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	synapseClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceKeyRead(d *pluginsdk.ResourceData, meta any) error {
	client := synapseClient.From(meta).KeysClient

	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.WorkspaceKeysID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.KeyName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Synapse Workspace Key %q (Workspace %q): %+v", id.KeyName, id.WorkspaceName, err)
	}

	workspaceID := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

	// Set the properties
	d.Set("synapse_workspace_id", workspaceID.ID())
	d.Set("active", resp.IsActiveCMK)
	d.Set("customer_managed_key_name", id.KeyName)
	d.Set("customer_managed_key_versionless_id", resp.KeyVaultURL)

	return nil
}

func synapseKeysWaitForStateChange(ctx context.Context, meta any, timeout time.Duration, resourceGroup string, workspaceName string, keyName string, pendingState string, targetState string) *pluginsdk.StateChangeConf {
	return &pluginsdk.StateChangeConf{
		Pending:    []string{pendingState},
		Target:     []string{targetState},
		MinTimeout: 1 * time.Minute,
		Timeout:    timeout,
		Refresh:    synapseKeysRefresh(ctx, meta, resourceGroup, workspaceName, keyName),
	}
}

func synapseKeysRefresh(ctx context.Context, meta any, resourceGroup string, workspaceName string, keyName string) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		client := synapseClient.From(meta).KeysClient

		log.Printf("[INFO] checking on state of encryption key '%q' (Workspace %q)", workspaceName, keyName)

		resp, err := client.Get(ctx, resourceGroup, workspaceName, keyName)
		if err != nil {
			return nil, "nil", fmt.Errorf("polling for the status of encryption key '%q' (Workspace %q): %v", keyName, workspaceName, err)
		}

		if resp.KeyProperties != nil && resp.IsActiveCMK != nil {
			return resp, strconv.FormatBool(*resp.IsActiveCMK), nil
		}

		// I am not returning an error here as this might have just been a bad get
		return resp, "nil", nil
	}
}
