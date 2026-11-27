// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_key

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceKeysDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.KeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceKeysID(d.Id())
	if err != nil {
		return err
	}

	// Fetch the key and check if it's an active key
	keyresult, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.KeyName)
	if err != nil {
		return fmt.Errorf("unable to fetch key %s in workspace %s: %v", id.KeyName, id.WorkspaceName, err)
	}

	// Azure only lets you delete keys that are not active
	if !*keyresult.IsActiveCMK {
		if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.KeyName); err != nil {
			return fmt.Errorf("unable to delete key %s in workspace %s: %v", id.KeyName, id.WorkspaceName, err)
		}
	}

	return nil
}
