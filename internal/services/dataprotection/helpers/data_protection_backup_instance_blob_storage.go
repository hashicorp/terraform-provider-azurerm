// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func PolicyProtectionStateRefreshFunc(ctx context.Context, client *backupinstanceresources.BackupInstanceResourcesClient, id backupinstanceresources.BackupInstanceId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.BackupInstancesGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
		}
		if res.Model == nil || res.Model.Properties == nil || res.Model.Properties.ProtectionStatus == nil || res.Model.Properties.ProtectionStatus.Status == nil {
			return nil, "", fmt.Errorf("reading DataProtection BackupInstance (%q) protection status: %+v", id, err)
		}

		return res, string(*res.Model.Properties.ProtectionStatus.Status), nil
	}
}
