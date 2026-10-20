// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_blob_storage

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupPolicyBlobStorageDelete(d *schema.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := basebackuppolicyresources.ParseBackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.BackupPoliciesDelete(ctx, *id); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting DataProtection BackupPolicy (%q): %+v", id, err)
	}
	return nil
}
