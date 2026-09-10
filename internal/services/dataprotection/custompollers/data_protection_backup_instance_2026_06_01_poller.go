// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &dataProtectionBackupInstance20260601Poller{}

type dataProtectionBackupInstance20260601Poller struct {
	client        *backupinstanceresources.BackupInstanceResourcesClient
	id            backupinstanceresources.BackupInstanceId
	pendingStates []backupinstanceresources.CurrentProtectionState
	targetState   backupinstanceresources.CurrentProtectionState
}

func NewDataProtectionBackupInstance20260601Poller(client *backupinstanceresources.BackupInstanceResourcesClient, id backupinstanceresources.BackupInstanceId, targetState backupinstanceresources.CurrentProtectionState, pendingStates []backupinstanceresources.CurrentProtectionState) *dataProtectionBackupInstance20260601Poller {
	return &dataProtectionBackupInstance20260601Poller{
		client:        client,
		id:            id,
		pendingStates: pendingStates,
		targetState:   targetState,
	}
}

func (p dataProtectionBackupInstance20260601Poller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.BackupInstancesGet(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", p.id)
	}

	if resp.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", p.id)
	}

	currentState := pointer.From(resp.Model.Properties.CurrentProtectionState)
	if currentState == p.targetState {
		return &pollers.PollResult{
			PollInterval: time.Minute,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	if slices.Contains(p.pendingStates, currentState) {
		return &pollers.PollResult{
			PollInterval: time.Minute,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	return nil, pollers.PollingFailedError{
		HttpResponse: &client.Response{Response: resp.HttpResponse},
		Message:      fmt.Sprintf("waiting for %s to reach state `%s` but got unexpected state `%s`", p.id, string(p.targetState), string(currentState)),
	}
}
