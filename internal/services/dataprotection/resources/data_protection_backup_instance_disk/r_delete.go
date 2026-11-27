// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_disk

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	dataprotectionClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupInstanceDiskDelete(d *schema.ResourceData, meta any) error {
	client := dataprotectionClient.From(meta).BackupInstanceClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := backupinstanceresources.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	if err = client.BackupInstancesDeleteThenPoll(ctx, *id, backupinstanceresources.DefaultBackupInstancesDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting DataProtection BackupInstance (%q): %+v", id, err)
	}

	return nil
}
