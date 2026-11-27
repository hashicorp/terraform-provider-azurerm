// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_encryption_scope

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/encryptionscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageEncryptionScopeDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.ResourceManager.EncryptionScopes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
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
