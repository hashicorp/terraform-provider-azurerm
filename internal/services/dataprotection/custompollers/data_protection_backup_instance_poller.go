// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &dataProtectionBackupInstancePoller{}

type dataProtectionBackupInstancePoller struct {
	client *backupinstanceresources.BackupInstanceResourcesClient
	id     backupinstanceresources.BackupInstanceId
}

var (
	backupInstancePollingSuccess = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	backupInstancePollingInProgress = pollers.PollResult{
		HttpResponse: nil,
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewDataProtectionBackupInstancePoller(client *backupinstanceresources.BackupInstanceResourcesClient, id backupinstanceresources.BackupInstanceId) *dataProtectionBackupInstancePoller {
	return &dataProtectionBackupInstancePoller{
		client: client,
		id:     id,
	}
}

func (p dataProtectionBackupInstancePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.BackupInstancesGet(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return nil, fmt.Errorf("polling for %s: `model` or `properties` was nil", p.id)
	}

	protectionState := pointer.From(resp.Model.Properties.CurrentProtectionState)
	if protectionState == backupinstanceresources.CurrentProtectionStateProtectionConfigured {
		return &backupInstancePollingSuccess, nil
	}

	return &backupInstancePollingInProgress, nil
}
