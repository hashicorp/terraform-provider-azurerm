// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2025-10-01-preview/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
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
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
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

	locks.ByName(id.WorkspaceName, "azurerm_databricks_workspace")
	defer locks.UnlockByName(id.WorkspaceName, "azurerm_databricks_workspace")

	existing, err := workspaceClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `Model` was nil", id)
	}

	if existing.Model.Properties.Parameters == nil {
		return fmt.Errorf("retrieving %s: `Parameters` was nil", id)
	}
	params := existing.Model.Properties.Parameters

	var encryptionEnabled bool
	if prepEncryption := params.PrepareEncryption; prepEncryption != nil {
		encryptionEnabled = prepEncryption.Value
	}

	if !encryptionEnabled {
		// TODO: consider removing this check and simply enabling this if it's not already?
		return fmt.Errorf("%s: `customer_managed_key_enabled` must be set to `true`", *id)
	}

	if params.Encryption != nil && params.Encryption.Value != nil && pointer.From(params.Encryption.Value.KeySource) != workspaces.KeySourceDefault {
		return tf.ImportAsExistsError("azurerm_databricks_workspace_root_dbfs_customer_managed_key", id.ID())
	}

	key, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}

	dbfsSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)
	keyVaultId := d.Get("key_vault_id").(string)
	if keyVaultId != "" {
		parsedKeyVaultID, err := commonids.ParseKeyVaultID(keyVaultId)
		if err != nil {
			return err
		}
		dbfsSubscriptionId = commonids.NewSubscriptionID(parsedKeyVaultID.SubscriptionId)
	}

	// make sure the key vault exists
	// TODO: consider removing this check and deprecating the `key_vault_id` property.
	// 1. The check can be time consuming when there are many KVs present in a subscription.
	// 2. The API request will fail if the KV doesn't exist, both that error and this check happen at apply time after a successful plan regardless.
	// 3. It doesn't work with Managed HSM vaults (hence the conditional)
	if !key.IsManagedHSM() {
		if keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, dbfsSubscriptionId, key.KeyVaultBaseURL); err != nil || keyVaultIdRaw == nil {
			return fmt.Errorf("retrieving the Resource ID for the Key Vault at URL %q: %+v", key.KeyVaultBaseURL, err)
		}
	}

	params.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource:   pointer.To(workspaces.KeySourceMicrosoftPointKeyvault),
			KeyName:     pointer.To(key.Name),
			Keyversion:  pointer.To(key.Version),
			Keyvaulturi: pointer.To(key.KeyVaultBaseURL),
		},
	}

	if err = workspaceClient.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
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
		if params := model.Properties.Parameters; params != nil {
			if encryption := params.Encryption; encryption != nil {
				if value := encryption.Value; value != nil {
					if strings.EqualFold(string(pointer.From(value.KeySource)), string(workspaces.KeySourceDefault)) && value.Keyvaulturi == nil && value.KeyName == nil {
						d.SetId("")
						return nil
					}

					key, err := keyvault.NewNestedItemID(pointer.From(value.Keyvaulturi), keyvault.NestedItemTypeKey, pointer.From(value.KeyName), pointer.From(value.Keyversion))
					if err != nil {
						return err
					}
					d.Set("key_vault_key_id", key.ID())
				}
			}
		}
	}

	d.Set("workspace_id", id.ID())
	d.Set("key_vault_id", d.Get("key_vault_id").(string))

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

	locks.ByName(id.WorkspaceName, "azurerm_databricks_workspace")
	defer locks.UnlockByName(id.WorkspaceName, "azurerm_databricks_workspace")

	existing, err := workspaceClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `Model` was nil", id)
	}

	if existing.Model.Properties.Parameters == nil {
		return fmt.Errorf("retrieving %s: `Parameters` was nil", id)
	}

	var encryptionEnabled bool
	if prepEncryption := existing.Model.Properties.Parameters.PrepareEncryption; prepEncryption != nil {
		encryptionEnabled = prepEncryption.Value
	}

	if !encryptionEnabled {
		return fmt.Errorf("%s: `customer_managed_key_enabled` must be set to `true`", *id)
	}

	key, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}

	dbfsSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)
	keyVaultId := d.Get("key_vault_id").(string)
	if keyVaultId != "" {
		parsedKeyVaultID, err := commonids.ParseKeyVaultID(keyVaultId)
		if err != nil {
			return err
		}
		dbfsSubscriptionId = commonids.NewSubscriptionID(parsedKeyVaultID.SubscriptionId)
	}

	// make sure the key vault exists
	if !key.IsManagedHSM() {
		if _, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, dbfsSubscriptionId, key.KeyVaultBaseURL); err != nil {
			return fmt.Errorf("retrieving the Resource ID for the Key Vault in subscription %q at URL %q: %+v", dbfsSubscriptionId, key.KeyVaultBaseURL, err)
		}
	}

	existing.Model.Properties.Parameters.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource:   pointer.To(workspaces.KeySourceMicrosoftPointKeyvault),
			KeyName:     pointer.To(key.Name),
			Keyversion:  pointer.To(key.Version),
			Keyvaulturi: pointer.To(key.KeyVaultBaseURL),
		},
	}

	if err = workspaceClient.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating Root DBFS Customer Managed Key for %s: %+v", *id, err)
	}

	// Always set this even if it's empty to keep the state file
	// consistent with the configuration file...
	d.Set("key_vault_id", keyVaultId)

	return databricksWorkspaceRootDbfsCustomerManagedKeyRead(d, meta)
}

func databricksWorkspaceRootDbfsCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.WorkspaceName, "azurerm_databricks_workspace")
	defer locks.UnlockByName(id.WorkspaceName, "azurerm_databricks_workspace")

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `Model` was nil", id)
	}

	if existing.Model.Properties.Parameters == nil {
		return fmt.Errorf("retrieving %s: `Parameters` was nil", id)
	}

	existing.Model.Properties.Parameters.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource: pointer.To(workspaces.KeySourceDefault),
		},
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("removing Root DBFS Customer Managed Key from %s: %+v", *id, err)
	}

	return nil
}
