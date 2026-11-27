// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_encryption_scope

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/encryptionscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageEncryptionScopeDelete(d *pluginsdk.ResourceData, meta any) error {
	client := storageClient.From(meta).ResourceManager.EncryptionScopes
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := encryptionscopes.ParseEncryptionScopeID(d.Id())
	if err != nil {
		return err
	}

	payload := encryptionscopes.EncryptionScope{
		Properties: &encryptionscopes.EncryptionScopeProperties{
			State: pointer.To(encryptionscopes.EncryptionScopeStateDisabled),
		},
	}

	if _, err = client.Put(ctx, *id, payload); err != nil {
		return fmt.Errorf("disabling %s: %+v", id, err)
	}

	return nil
}
