// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_workspace_root_dbfs_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func databricksWorkspaceRootDbfsCustomerManagedKeyCreate(d *pluginsdk.ResourceData, meta any) error {
	workspaceClient := meta.(*clients.Client).DataBricks.WorkspacesClient

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

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		if params.Encryption != nil && params.Encryption.Value != nil && pointer.From(params.Encryption.Value.KeySource) != workspaces.KeySourceDefault {
			return tf.ImportAsExistsError("azurerm_databricks_workspace_root_dbfs_customer_managed_key", id.ID())
		}
	}

	key, err := keyvault.ParseNestedItemID(d.Get("key_vault_key_id").(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}

	params.Encryption = &workspaces.WorkspaceEncryptionParameter{
		Value: &workspaces.Encryption{
			KeySource:   pointer.To(workspaces.KeySourceMicrosoftPointKeyvault),
			KeyName:     pointer.To(key.Name),
			Keyversion:  pointer.To(key.Version),
			Keyvaulturi: pointer.To(key.KeyVaultBaseURL),
		},
	}

	if err = workspaceClient.CreateOrUpdateCallbackThenPoll(ctx, *id, *existing.Model, sdk.SetIDAndIdentityCallback(meta, id, d)); err != nil {
		return fmt.Errorf("creating Root DBFS Customer Managed Key for %s: %+v", *id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, id); err != nil {
		return err
	}

	return databricksWorkspaceRootDbfsCustomerManagedKeyRead(d, meta)
}
