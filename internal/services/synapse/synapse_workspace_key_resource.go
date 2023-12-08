// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseWorkspaceKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseWorkspaceKeysCreateUpdate,
		Read:   resourceSynapseWorkspaceKeyRead,
		Update: resourceSynapseWorkspaceKeysCreateUpdate,
		Delete: resourceSynapseWorkspaceKeysDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceKeysID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"customer_managed_key_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"customer_managed_key_versionless_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.VersionlessNestedItemId,
			},

			"active": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSynapseWorkspaceKeysCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.KeysClient
	// workspaceClient := meta.(*clients.Client).Synapse.WorkspaceClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	key := d.Get("customer_managed_key_versionless_id")
	keyName := d.Get("customer_managed_key_name").(string)
	isActiveCMK := d.Get("active").(bool)

	log.Printf("[INFO] Is active CMK: %t", isActiveCMK)

	keyProperties := synapse.KeyProperties{
		IsActiveCMK: &isActiveCMK,
		KeyVaultURL: utils.String(key.(string)),
	}

	synapseKey := synapse.Key{
		KeyProperties: &keyProperties,
	}

	actualKeyName := ""
	if keyName != "" {
		actualKeyName = keyName
	}

	locks.ByName(workspaceId.Name, "azurerm_synapse_workspace")
	defer locks.UnlockByName(workspaceId.Name, "azurerm_synapse_workspace")
	keyresult, err := client.CreateOrUpdate(ctx, workspaceId.ResourceGroup, workspaceId.Name, actualKeyName, synapseKey)
	if err != nil {
		return fmt.Errorf("creating Synapse Workspace Key %q (Workspace %q): %+v", workspaceId.Name, workspaceId.Name, err)
	}

	if keyresult.ID == nil || *keyresult.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse Key 'cmk'")
	}

	// If the state of the key in the response (from Azure) is not equal to the desired target state (from plan/config), we'll wait until that change is complete
	if isActiveCMK != *keyresult.KeyProperties.IsActiveCMK {
		updateWait := synapseKeysWaitForStateChange(ctx, meta, d.Timeout(pluginsdk.TimeoutUpdate), workspaceId.ResourceGroup, workspaceId.Name, actualKeyName, strconv.FormatBool(*keyresult.KeyProperties.IsActiveCMK), strconv.FormatBool(isActiveCMK))

		if _, err := updateWait.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for Synapse Keys to finish updating '%q' (Workspace Group %q): %v", actualKeyName, workspaceId.Name, err)
		}
	}

	id := parse.NewWorkspaceKeysID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, actualKeyName)
	d.SetId(id.ID())

	return resourceSynapseWorkspaceKeyRead(d, meta)
}

func resourceSynapseWorkspaceKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.KeysClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceKeysID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.KeyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Synapse Workspace Key %q (Workspace %q): %+v", id.KeyName, id.WorkspaceName, err)
	}

	workspaceID := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

	// Set the properties
	d.Set("synapse_workspace_id", workspaceID.ID())
	d.Set("active", resp.KeyProperties.IsActiveCMK)
	d.Set("customer_managed_key_name", id.KeyName)
	d.Set("customer_managed_key_versionless_id", resp.KeyProperties.KeyVaultURL)

	return nil
}

func resourceSynapseWorkspaceKeysDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
	if !*keyresult.KeyProperties.IsActiveCMK {
		_, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.KeyName)
		if err != nil {
			return fmt.Errorf("unable to delete key %s in workspace %s: %v", id.KeyName, id.WorkspaceName, err)
		}
	}

	return nil
}

func synapseKeysWaitForStateChange(ctx context.Context, meta interface{}, timeout time.Duration, resourceGroup string, workspaceName string, keyName string, pendingState string, targetState string) *pluginsdk.StateChangeConf {
	return &pluginsdk.StateChangeConf{
		Pending:    []string{pendingState},
		Target:     []string{targetState},
		MinTimeout: 1 * time.Minute,
		Timeout:    timeout,
		Refresh:    synapseKeysRefresh(ctx, meta, resourceGroup, workspaceName, keyName),
	}
}

func synapseKeysRefresh(ctx context.Context, meta interface{}, resourceGroup string, workspaceName string, keyName string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*clients.Client).Synapse.KeysClient

		log.Printf("[INFO] checking on state of encryption key '%q' (Workspace %q)", workspaceName, keyName)

		resp, err := client.Get(ctx, resourceGroup, workspaceName, keyName)
		if err != nil {
			return nil, "nil", fmt.Errorf("polling for the status of encryption key '%q' (Workspace %q): %v", keyName, workspaceName, err)
		}

		if resp.KeyProperties != nil && resp.KeyProperties.IsActiveCMK != nil {
			return resp, strconv.FormatBool(*resp.KeyProperties.IsActiveCMK), nil
		}

		// I am not returning an error here as this might have just been a bad get
		return resp, "nil", nil
	}
}
