// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompoller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type snapshotCopyStartPoller struct {
	client *snapshots.SnapshotsClient
	id     snapshots.SnapshotId
}

var _ pollers.PollerType = &snapshotCopyStartPoller{}

func NewSnapshotCopyStartPoller(client *snapshots.SnapshotsClient, id snapshots.SnapshotId) *snapshotCopyStartPoller {
	return &snapshotCopyStartPoller{
		client: client,
		id:     id,
	}
}

func (s snapshotCopyStartPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := s.client.Get(ctx, s.id)
	if err != nil {
		return nil, fmt.Errorf("polling for the copy status of %s: %+v", s.id, err)
	}

	completionPercent := float64(0)
	if model := resp.Model; model != nil && model.Properties != nil {
		completionPercent = pointer.From(model.Properties.CompletionPercent)
	}

	log.Printf("[DEBUG] Snapshot %s copy completion is at %.1f%%", s.id, completionPercent)
	if completionPercent < 100 {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusInProgress,
			PollInterval: 30 * time.Second,
		}, nil
	}

	return &pollers.PollResult{
		Status: pollers.PollingStatusSucceeded,
	}, nil
}
