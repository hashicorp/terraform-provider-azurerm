// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_transparent_data_encryption

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceencryptionprotectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstancekeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlManagedInstanceTransparentDataEncryptionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceEncryptionProtectorClient
	managedInstanceKeysClient := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceKeysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managedInstanceId, err := commonids.ParseSqlManagedInstanceID(d.Get("managed_instance_id").(string))
	if err != nil {
		return err
	}

	// Encryption protector always uses "current" for the name
	id := parse.NewManagedInstanceEncryptionProtectorID(managedInstanceId.SubscriptionId, managedInstanceId.ResourceGroupName, managedInstanceId.ManagedInstanceName, "current")

	// Normally we would check if this is a new resource, but the way encryption protector works, it always overwrites
	// whatever is there anyways. Compounding the issue is that SQL Server creates an instance of encryption protector
	// which causes the isNewResource check to fail because we are trying to create the encryption as a separate resource
	// and encryption protector is already present. The reason we create encryption protector as a separate resource is
	// because after the SQL server is created, we need to grant it permissions to AKV, so encryption protector can use those
	// keys are part of setting up TDE

	payload := managedinstanceencryptionprotectors.ManagedInstanceEncryptionProtector{
		Properties: &managedinstanceencryptionprotectors.ManagedInstanceEncryptionProtectorProperties{
			AutoRotationEnabled: pointer.To(d.Get("auto_rotation_enabled").(bool)),
			ServerKeyName:       pointer.To(""),
			ServerKeyType:       managedinstanceencryptionprotectors.ServerKeyTypeServiceManaged,
		},
	}

	var key *keyvault.NestedItemID
	if v, ok := d.GetOk("key_vault_key_id"); ok {
		key, err = keyvault.ParseNestedItemID(v.(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
		if err != nil {
			return err
		}

		key, err = resourceMsSqlManagedInstanceTransparentDataEncryptionVersionedKey(ctx, key, "key_vault_key_id", d.Get("auto_rotation_enabled").(bool), meta)
		if err != nil {
			return err
		}
	}

	if key != nil {
		keyVaultName, err := resourceMsSqlManagedInstanceTransparentDataEncryptionKeyVaultName(key.KeyVaultBaseURL)
		if err != nil {
			return err
		}

		// Key name value for the Encryption Protector. Format is: {vaultName}_{key}_{key_version}
		managedInstanceKeyName := fmt.Sprintf("%s_%s_%s", keyVaultName, key.Name, key.Version)

		managedInstanceKeyId := managedinstancekeys.NewManagedInstanceKeyID(managedInstanceId.SubscriptionId, managedInstanceId.ResourceGroupName, managedInstanceId.ManagedInstanceName, managedInstanceKeyName)
		managedInstanceKeyPayload := managedinstancekeys.ManagedInstanceKey{
			Properties: &managedinstancekeys.ManagedInstanceKeyProperties{
				AutoRotationEnabled: pointer.To(d.Get("auto_rotation_enabled").(bool)),
				ServerKeyType:       managedinstancekeys.ServerKeyTypeAzureKeyVault,
				Uri:                 pointer.To(key.ID()),
			},
		}

		if err := managedInstanceKeysClient.CreateOrUpdateThenPoll(ctx, managedInstanceKeyId, managedInstanceKeyPayload); err != nil {
			return fmt.Errorf("creating %s: %+v", managedInstanceKeyId, err)
		}

		// Update TDE properties to reflect usage of Managed Key
		payload.Properties.ServerKeyName = pointer.To(managedInstanceKeyName)
		payload.Properties.ServerKeyType = managedinstanceencryptionprotectors.ServerKeyTypeAzureKeyVault
	}

	// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
	if err := client.CreateOrUpdateThenPoll(ctx, *managedInstanceId, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMsSqlManagedInstanceTransparentDataEncryptionRead(d, meta)
}
