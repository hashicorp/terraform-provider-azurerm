// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2024-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksWorkspaceRootDbfsCustomerManagedKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: databricksWorkspaceRootDbfsCustomerManagedKeyCreate,
		Read:   databricksWorkspaceRootDbfsCustomerManagedKeyRead,
		Update: databricksWorkspaceRootDbfsCustomerManagedKeyUpdate,
		Delete: databricksWorkspaceRootDbfsCustomerManagedKeyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := workspaces.ParseWorkspaceID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			// validate that the passed ID is a valid CMK configuration ID
			id, err := workspaces.ParseWorkspaceID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, err
			}

			// set the new values for the CMK resource
			d.SetId(id.ID())
			d.Set("workspace_id", id.ID())

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: map[string]*pluginsdk.Schema{
			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: keyVaultValidate.KeyVaultChildID,
			},

			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateKeyVaultID,
			},
		},
	}
}

func databricksWorkspaceRootDbfsCustomerManagedKeyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	workspaceClient := meta.(*clients.Client).DataBricks.WorkspacesClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}

	var keyIdRaw string
	var keyVaultId string
	var key *keyVaultParse.NestedItemId

	if v, ok := d.GetOk("key_vault_key_id"); ok {
		keyIdRaw = v.(string)
	}

	if v, ok := d.GetOk("key_vault_id"); ok {
		keyVaultId = v.(string)
	}

	key, err = keyVaultParse.ParseNestedItemID(keyIdRaw)
	if err != nil {
		return err
	}

	// Not sure if I should also lock the key vault here too
	// or at the very least the key?
	locks.ByName(id.WorkspaceName, "azurerm_databricks_workspace")
	defer locks.UnlockByName(id.WorkspaceName, "azurerm_databricks_workspace")
	var encryptionEnabled bool

	workspace, err := workspaceClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keySource := workspaces.KeySourceDefault
	var params *workspaces.WorkspaceCustomParameters

	if model := workspace.Model; model != nil {
		if params = model.Properties.Parameters; params != nil {
			if params.PrepareEncryption != nil {
				encryptionEnabled = model.Properties.Parameters.PrepareEncryption.Value
			}

			if params.Encryption != nil && params.Encryption.Value != nil && params.Encryption.Value.KeySource != nil {
				keySource = pointer.From(params.Encryption.Value.KeySource)
			}
		} else {
			return fmt.Errorf("`WorkspaceCustomParameters` were nil")
		}
	} else {
		return fmt.Errorf("`Workspace` was nil")
	}

	if !encryptionEnabled {
		return fmt.Errorf("%s: `customer_managed_key_enabled` must be set to `true`", *id)
	}

	// If the 'root_dbfs_cmk_key_vault_id' was not defined assume the
	// key vault exists in the same subscription as the workspace...
	dbfsSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	// If they passed the 'root_dbfs_cmk_key_vault_id' parse the Key Vault ID
	// to extract the correct key vault subscription for the exists call...
	if keyVaultId != "" {
		keyVaultId, err := commonids.ParseKeyVaultID(keyVaultId)
		if err != nil {
			return fmt.Errorf("parsing %q as a Key Vault ID: %+v", keyVaultId, err)
		}

		dbfsSubscriptionId = commonids.NewSubscriptionID(keyVaultId.SubscriptionId)
	}

	// make sure the key vault exists
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, dbfsSubscriptionId, key.KeyVaultBaseUrl)
	if err != nil || keyVaultIdRaw == nil {
		return fmt.Errorf("retrieving the Resource ID for the Key Vault at URL %q: %+v", key.KeyVaultBaseUrl, err)
	}

	// Only throw the import error if the keysource value has been set to something other than default...
	if params.Encryption != nil && params.Encryption.Value != nil && keySource != workspaces.KeySourceDefault {
		return tf.ImportAsExistsError("azurerm_databricks_workspace_root_dbfs_customer_managed_key", id.ID())
	}

	// We need to pull all of the custom params from the parent
	// workspace resource and then add our new encryption values into the
	// structure, else the other values set in the parent workspace
	// resource will be lost and overwritten as nil. ¯\_(ツ)_/¯
	// NOTE: 'workspace.Parameters' will never be nil as 'customer_managed_key_enabled' and 'infrastructure_encryption_enabled'
	// fields have a default value in the parent workspace resource.
	params.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource:   pointer.To(workspaces.KeySourceMicrosoftPointKeyvault),
			KeyName:     pointer.To(key.Name),
			Keyversion:  pointer.To(key.Version),
			Keyvaulturi: pointer.To(key.KeyVaultBaseUrl),
		},
	}

	props := pointer.From(workspace.Model)
	props.Properties.Parameters = params

	if err = workspaceClient.CreateOrUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("creating Root DBFS Customer Managed Key for %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	// Always set this even if it's empty to keep the state file
	// consistent with the configuration file...
	d.Set("key_vault_id", keyVaultId)

	return databricksWorkspaceRootDbfsCustomerManagedKeyRead(d, meta)
}

func databricksWorkspaceRootDbfsCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	var keyVaultId string
	if v, ok := d.GetOk("key_vault_id"); ok {
		keyVaultId = v.(string)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if model.Properties.Parameters != nil {
			if props := model.Properties.Parameters.Encryption; props != nil {
				if strings.EqualFold(string(*props.Value.KeySource), string(workspaces.KeySourceMicrosoftPointKeyvault)) && (props.Value.KeyName == nil || props.Value.Keyversion == nil || props.Value.Keyvaulturi == nil) {
					d.SetId("")
					return nil
				}

				key, err := keyVaultParse.NewNestedItemID(*props.Value.Keyvaulturi, keyVaultParse.NestedItemTypeKey, *props.Value.KeyName, *props.Value.Keyversion)
				if err == nil {
					d.Set("key_vault_key_id", key.ID())
				}
			}
		}
	}

	d.Set("workspace_id", id.ID())
	d.Set("key_vault_id", keyVaultId)

	return nil
}

func databricksWorkspaceRootDbfsCustomerManagedKeyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	workspaceClient := meta.(*clients.Client).DataBricks.WorkspacesClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}

	var key *keyVaultParse.NestedItemId
	var params *workspaces.WorkspaceCustomParameters
	var keyVaultId string
	var keyVaultKeyId string

	if v, ok := d.GetOk("key_vault_key_id"); ok {
		keyVaultKeyId = v.(string)
	}

	if v, ok := d.GetOk("key_vault_id"); ok {
		keyVaultId = v.(string)
	}

	key, err = keyVaultParse.ParseNestedItemID(keyVaultKeyId)
	if err != nil {
		return err
	}

	// Not sure if I should also lock the key vault here too
	// or at the very least the key?
	locks.ByName(id.WorkspaceName, "azurerm_databricks_workspace")
	defer locks.UnlockByName(id.WorkspaceName, "azurerm_databricks_workspace")
	var encryptionEnabled bool

	workspace, err := workspaceClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := workspace.Model; model != nil {
		if params = model.Properties.Parameters; params != nil {
			if params.PrepareEncryption != nil {
				encryptionEnabled = model.Properties.Parameters.PrepareEncryption.Value
			}
		} else {
			return fmt.Errorf("`WorkspaceCustomParameters` were nil")
		}
	} else {
		return fmt.Errorf("`Workspace` was nil")
	}

	if !encryptionEnabled {
		return fmt.Errorf("%s: `customer_managed_key_enabled` must be set to `true`", *id)
	}

	// If the 'root_dbfs_cmk_key_vault_id' was not defined assume
	// the key vault exists in the same subscription as the workspace...
	dbfsSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	if keyVaultId != "" {
		v, err := commonids.ParseKeyVaultID(keyVaultId)
		if err != nil {
			return fmt.Errorf("parsing %q as a Key Vault ID: %+v", keyVaultId, err)
		}

		dbfsSubscriptionId = commonids.NewSubscriptionID(v.SubscriptionId)
	}

	// make sure the key vault exists
	_, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, dbfsSubscriptionId, key.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID for the Key Vault in subscription %q at URL %q: %+v", dbfsSubscriptionId, key.KeyVaultBaseUrl, err)
	}

	// We need to pull all of the custom params from the parent
	// workspace resource and then add our new encryption values into the
	// structure, else the other values set in the parent workspace
	// resource will be lost and overwritten as nil. ¯\_(ツ)_/¯
	// NOTE: 'workspace.Parameters' will never be nil as 'customer_managed_key_enabled' and 'infrastructure_encryption_enabled'
	// fields have a default value in the parent workspace resource.
	params.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource:   pointer.To(workspaces.KeySourceMicrosoftPointKeyvault),
			KeyName:     pointer.To(key.Name),
			Keyversion:  pointer.To(key.Version),
			Keyvaulturi: pointer.To(key.KeyVaultBaseUrl),
		},
	}

	props := pointer.From(workspace.Model)
	props.Properties.Parameters = params

	if err = workspaceClient.CreateOrUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating Root DBFS Customer Managed Key for %s: %+v", *id, err)
	}

	// Always set this even if it's empty to keep the state file
	// consistent with the configuration file...
	d.Set("key_vault_id", keyVaultId)

	return databricksWorkspaceRootDbfsCustomerManagedKeyRead(d, meta)
}

func databricksWorkspaceRootDbfsCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	// Not sure if I should also lock the key vault here too
	locks.ByName(id.WorkspaceName, "azurerm_databricks_workspace")
	defer locks.UnlockByName(id.WorkspaceName, "azurerm_databricks_workspace")

	workspace, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if workspace.Model == nil {
		return fmt.Errorf("`Workspace` was nil")
	}

	if workspace.Model.Properties.Parameters == nil {
		return fmt.Errorf("`WorkspaceCustomParameters` were nil")
	}

	// Since this isn't real and you cannot turn off CMK without destroying the
	// workspace and recreating it the best I can do is to set the workspace
	// back to using Microsoft managed keys and removing the CMK fields
	// also need to pull all of the custom params from the parent
	// workspace resource and then add our new encryption values into the
	// structure, else the other values set in the parent workspace
	// resource will be lost and overwritten as nil. ¯\_(ツ)_/¯
	params := workspace.Model.Properties.Parameters
	params.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource: pointer.To(workspaces.KeySourceDefault),
		},
	}

	props := pointer.From(workspace.Model)
	props.Properties.Parameters = params

	if err = client.CreateOrUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("removing Root DBFS Customer Managed Key from %s: %+v", *id, err)
	}

	return nil
}
