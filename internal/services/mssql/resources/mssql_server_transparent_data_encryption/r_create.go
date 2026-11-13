// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_transparent_data_encryption

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/encryptionprotectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serverkeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlTransparentDataEncryptionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	serverKeysClient := meta.(*clients.Client).MSSQL.ServerKeysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	// Encryption protector always uses "current" for the name
	id := parse.NewEncryptionProtectorID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, "current")

	// Normally we would check if this is a new resource, but the way encryption protector works, it always overwrites
	// whatever is there anyways. Compounding the issue is that SQL Server creates an instance of encryption protector
	// which causes the isNewResource check to fail because we are trying to create the encryption as a separate resource
	// and encryption protector is already present. The reason we create encryption protector as a separate resource is
	// because after the SQL server is created, we need to grant it permissions to AKV, so encryption protector can use those
	// keys are part of setting up TDE

	payload := encryptionprotectors.EncryptionProtector{
		Properties: &encryptionprotectors.EncryptionProtectorProperties{
			AutoRotationEnabled: pointer.To(d.Get("auto_rotation_enabled").(bool)),
			ServerKeyName:       pointer.To(""),
			ServerKeyType:       encryptionprotectors.ServerKeyTypeServiceManaged,
		},
	}

	var key *keyvault.NestedItemID
	if v, ok := d.GetOk("key_vault_key_id"); ok {
		keyId, err := keyvault.ParseNestedItemID(v.(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return err
		}
		key = keyId
	}

	if key != nil {
		keyVaultName, err := resourceMsSqlTransparentDataEncryptionKeyVaultName(key.KeyVaultBaseURL)
		if err != nil {
			return err
		}

		serverKeyName := fmt.Sprintf("%s_%s_%s", keyVaultName, key.Name, key.Version)

		serverKeyId := serverkeys.NewKeyID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, serverKeyName)
		serverKeyPayload := serverkeys.ServerKey{
			Properties: &serverkeys.ServerKeyProperties{
				AutoRotationEnabled: pointer.To(d.Get("auto_rotation_enabled").(bool)),
				ServerKeyType:       serverkeys.ServerKeyTypeAzureKeyVault,
				Uri:                 pointer.To(key.ID()),
			},
		}

		if err := serverKeysClient.CreateOrUpdateThenPoll(ctx, serverKeyId, serverKeyPayload); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", serverKeyId, err)
		}

		// Update TDE properties to reflect usage of Server Key
		payload.Properties.ServerKeyName = pointer.To(serverKeyName)
		payload.Properties.ServerKeyType = encryptionprotectors.ServerKeyTypeAzureKeyVault
	}

	// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
	if err := client.CreateOrUpdateThenPoll(ctx, *serverId, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceMsSqlTransparentDataEncryptionRead(d, meta)
}
