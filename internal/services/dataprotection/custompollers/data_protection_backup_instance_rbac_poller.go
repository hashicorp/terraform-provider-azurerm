// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &dataProtectionBackupInstanceRBACPoller{}

type dataProtectionBackupInstanceRBACPoller struct {
	client     *backupvaultresources.BackupVaultResourcesClient
	id         backupvaultresources.BackupVaultId
	parameters backupvaultresources.ValidateForBackupRequest
}

func NewDataProtectionBackupInstanceRBACPoller(client *backupvaultresources.BackupVaultResourcesClient, id backupvaultresources.BackupVaultId, parameters backupvaultresources.ValidateForBackupRequest) *dataProtectionBackupInstanceRBACPoller {
	return &dataProtectionBackupInstanceRBACPoller{
		client:     client,
		id:         id,
		parameters: parameters,
	}
}

func (p dataProtectionBackupInstanceRBACPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if err := p.client.BackupInstancesValidateForBackupThenPoll(ctx, p.id, p.parameters); err != nil {
		if ctx.Err() != nil {
			return nil, err
		}

		if strings.Contains(err.Error(), "UserErrorExtensionMSIMissingPermissions") ||
			strings.Contains(err.Error(), "UserErrorMissingMSIPermissionsOnSnapshotResourceGroup") ||
			 strings.Contains(err.Error(), "UserErrorBackupStorageLocationNotReady") {
			log.Printf("[DEBUG] Waiting for backup RBAC permissions to propagate for %s: %+v", p.id, err)
			return &pollers.PollResult{
				PollInterval: 30 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}

		return nil, fmt.Errorf("validating backup permissions for %s: %+v", p.id, err)
	}

	return &pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
