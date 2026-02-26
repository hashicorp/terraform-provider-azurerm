// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &dataProtectionBackupVaultPoller{}

type dataProtectionBackupVaultPoller struct {
	client *backupvaultresources.BackupVaultResourcesClient
	id     backupvaultresources.BackupVaultId
}

var (
	pollingSuccess = pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	pollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewDataProtectionBackupVaultPoller(client *backupvaultresources.BackupVaultResourcesClient, id backupvaultresources.BackupVaultId) *dataProtectionBackupVaultPoller {
	return &dataProtectionBackupVaultPoller{
		client: client,
		id:     id,
	}
}

func (p dataProtectionBackupVaultPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.BackupVaultsGet(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollingSuccess, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &pollingInProgress, nil
}
