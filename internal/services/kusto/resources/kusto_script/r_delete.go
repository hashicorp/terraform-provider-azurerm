// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto_script

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/scripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKustoDatabaseScriptDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Kusto.ScriptsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scripts.ParseScriptID(d.Id())
	if err != nil {
		return err
	}

	// DELETE operation for script does not support running concurrently at cluster level
	locks.ByName(id.ClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(id.ClusterName, "azurerm_kusto_cluster")

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}
