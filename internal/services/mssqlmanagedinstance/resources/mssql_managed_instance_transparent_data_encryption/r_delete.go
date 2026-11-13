// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_transparent_data_encryption

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceencryptionprotectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlManagedInstanceTransparentDataEncryptionDelete(d *pluginsdk.ResourceData, meta any) error {
	// Note that encryption protector cannot be deleted. It can only be updated between AzureKeyVault
	// and SystemManaged. For safety, when this resource is deleted, we're resetting the key type
	// to service managed to prevent accidental lockout if someone were to delete the keys from key vault

	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceEncryptionProtectorClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstanceEncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	encryptionProtector := managedinstanceencryptionprotectors.ManagedInstanceEncryptionProtector{
		Properties: &managedinstanceencryptionprotectors.ManagedInstanceEncryptionProtectorProperties{
			ServerKeyType: managedinstanceencryptionprotectors.ServerKeyTypeServiceManaged,
			ServerKeyName: pointer.To(""),
		},
	}

	if err = client.CreateOrUpdateThenPoll(ctx, managedInstanceId, encryptionProtector); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
