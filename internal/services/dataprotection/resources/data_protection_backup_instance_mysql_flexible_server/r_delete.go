// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_mysql_flexible_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.BackupInstancesDeleteThenPoll(ctx, *id, backupinstanceresources.DefaultBackupInstancesDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
