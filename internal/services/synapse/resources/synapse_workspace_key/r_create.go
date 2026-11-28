// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_workspace_key

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceKeysCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Synapse.KeysClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	// TODO: import check?

	key := d.Get("customer_managed_key_versionless_id")
	keyName := d.Get("customer_managed_key_name").(string)
	isActiveCMK := d.Get("active").(bool)

	log.Printf("[INFO] Is active CMK: %t", isActiveCMK)

	keyProperties := synapse.KeyProperties{
		IsActiveCMK: &isActiveCMK,
		KeyVaultURL: pointer.To(key.(string)),
	}

	synapseKey := synapse.Key{
		KeyProperties: &keyProperties,
	}

	actualKeyName := ""
	if keyName != "" {
		actualKeyName = keyName
	}

	locks.ByName(workspaceId.WorkspaceName, "azurerm_synapse_workspace")
	defer locks.UnlockByName(workspaceId.WorkspaceName, "azurerm_synapse_workspace")
	keyresult, err := client.CreateOrUpdate(ctx, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, actualKeyName, synapseKey)
	if err != nil {
		return fmt.Errorf("creating Synapse Workspace Key %q (Workspace %q): %+v", workspaceId.WorkspaceName, workspaceId.WorkspaceName, err)
	}

	if keyresult.ID == nil || *keyresult.ID == "" {
		return errors.New("empty or nil ID returned for Synapse Key 'cmk'")
	}

	if d.IsNewResource() {
		id := parse.NewWorkspaceKeysID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, actualKeyName)
		d.SetId(id.ID())
	}

	// If the state of the key in the response (from Azure) is not equal to the desired target state (from plan/config), we'll wait until that change is complete
	if isActiveCMK != *keyresult.IsActiveCMK {
		updateWait := synapseKeysWaitForStateChange(ctx, meta, d.Timeout(pluginsdk.TimeoutUpdate), workspaceId.ResourceGroupName, workspaceId.WorkspaceName, actualKeyName, strconv.FormatBool(*keyresult.IsActiveCMK), strconv.FormatBool(isActiveCMK))

		if _, err := updateWait.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for Synapse Keys to finish updating '%q' (Workspace Group %q): %v", actualKeyName, workspaceId.WorkspaceName, err)
		}
	}

	return resourceSynapseWorkspaceKeyRead(d, meta)
}
