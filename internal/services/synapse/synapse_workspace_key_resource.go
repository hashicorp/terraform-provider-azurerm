// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/keys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSynapseWorkspaceKey() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
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
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeKey),
			},

			"active": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}

	if !features.FivePointOh() {
		r.Schema["customer_managed_key_versionless_id"].ValidateFunc = keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeAny)
	}

	return r
}

func resourceSynapseWorkspaceKeysCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.KeysClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	// TODO: import check?

	key := d.Get("customer_managed_key_versionless_id")
	keyName := d.Get("customer_managed_key_name").(string)
	isActiveCMK := d.Get("active").(bool)

	log.Printf("[INFO] Is active CMK: %t", isActiveCMK)

	actualKeyName := ""
	if keyName != "" {
		actualKeyName = keyName
	}

	id := keys.NewKeyID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, actualKeyName)

	synapseKey := keys.Key{
		Properties: &keys.KeyProperties{
			IsActiveCMK: pointer.To(isActiveCMK),
			KeyVaultURL: pointer.To(key.(string)),
		},
	}

	locks.ByName(workspaceId.Name, "azurerm_synapse_workspace")
	defer locks.UnlockByName(workspaceId.Name, "azurerm_synapse_workspace")
	keyresult, err := client.CreateOrUpdate(ctx, id, synapseKey)
	if err != nil {
		return fmt.Errorf("creating Synapse Workspace Key %q (Workspace %q): %+v", workspaceId.Name, workspaceId.Name, err)
	}

	if keyresult.Model == nil || keyresult.Model.Id == nil || *keyresult.Model.Id == "" {
		return fmt.Errorf("empty or nil ID returned for Synapse Key 'cmk'")
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	resultIsActiveCMK := false
	if keyresult.Model.Properties != nil && keyresult.Model.Properties.IsActiveCMK != nil {
		resultIsActiveCMK = *keyresult.Model.Properties.IsActiveCMK
	}

	// If the state of the key in the response (from Azure) is not equal to the desired target state (from plan/config), we'll wait until that change is complete
	if isActiveCMK != resultIsActiveCMK {
		updateWait := synapseKeysWaitForStateChange(ctx, meta, d.Timeout(pluginsdk.TimeoutUpdate), id, strconv.FormatBool(resultIsActiveCMK), strconv.FormatBool(isActiveCMK))

		if _, err := updateWait.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for Synapse Keys to finish updating '%q' (Workspace Group %q): %v", actualKeyName, workspaceId.Name, err)
		}
	}

	return resourceSynapseWorkspaceKeyRead(d, meta)
}

func resourceSynapseWorkspaceKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.KeysClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keys.ParseKeyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Synapse Workspace Key %q (Workspace %q): %+v", id.KeyName, id.WorkspaceName, err)
	}

	workspaceID := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

	// Set the properties
	d.Set("synapse_workspace_id", workspaceID.ID())
	d.Set("customer_managed_key_name", id.KeyName)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("active", props.IsActiveCMK)
			d.Set("customer_managed_key_versionless_id", props.KeyVaultURL)
		}
	}

	return nil
}

func resourceSynapseWorkspaceKeysDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.KeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keys.ParseKeyID(d.Id())
	if err != nil {
		return err
	}

	// Fetch the key and check if it's an active key
	keyresult, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("unable to fetch key %s in workspace %s: %v", id.KeyName, id.WorkspaceName, err)
	}

	// Azure only lets you delete keys that are not active
	isActiveCMK := false
	if keyresult.Model != nil && keyresult.Model.Properties != nil && keyresult.Model.Properties.IsActiveCMK != nil {
		isActiveCMK = *keyresult.Model.Properties.IsActiveCMK
	}
	if !isActiveCMK {
		if _, err := client.Delete(ctx, *id); err != nil {
			return fmt.Errorf("unable to delete key %s in workspace %s: %v", id.KeyName, id.WorkspaceName, err)
		}
	}

	return nil
}

func synapseKeysWaitForStateChange(ctx context.Context, meta interface{}, timeout time.Duration, id keys.KeyId, pendingState string, targetState string) *pluginsdk.StateChangeConf {
	return &pluginsdk.StateChangeConf{
		Pending:    []string{pendingState},
		Target:     []string{targetState},
		MinTimeout: 1 * time.Minute,
		Timeout:    timeout,
		Refresh:    synapseKeysRefresh(ctx, meta, id),
	}
}

func synapseKeysRefresh(ctx context.Context, meta interface{}, id keys.KeyId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*clients.Client).Synapse.KeysClient

		log.Printf("[INFO] checking on state of encryption key '%q' (Workspace %q)", id.WorkspaceName, id.KeyName)

		resp, err := client.Get(ctx, id)
		if err != nil {
			return nil, "nil", fmt.Errorf("polling for the status of encryption key '%q' (Workspace %q): %v", id.KeyName, id.WorkspaceName, err)
		}

		if model := resp.Model; model != nil && model.Properties != nil && model.Properties.IsActiveCMK != nil {
			return resp, strconv.FormatBool(*model.Properties.IsActiveCMK), nil
		}

		// I am not returning an error here as this might have just been a bad get
		return resp, "nil", nil
	}
}
