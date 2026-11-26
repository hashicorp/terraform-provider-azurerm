// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_transparent_data_encryption

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/encryptionprotectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlTransparentDataEncryptionDelete(d *pluginsdk.ResourceData, meta any) error {
	// Note that encryption protector cannot be deleted. It can only be updated between AzureKeyVault
	// and SystemManaged. For safety, when this resource is deleted, we're resetting the key type
	// to service managed to prevent accidental lockout if someone were to delete the keys from key vault

	client := meta.(*clients.Client).MSSQL.EncryptionProtectorClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	encryptionProtector := encryptionprotectors.EncryptionProtector{
		Properties: &encryptionprotectors.EncryptionProtectorProperties{
			ServerKeyType: encryptionprotectors.ServerKeyTypeServiceManaged,
			ServerKeyName: pointer.To(""),
		},
	}

	if err = client.CreateOrUpdateThenPoll(ctx, serverId, encryptionProtector); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
