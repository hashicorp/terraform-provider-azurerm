// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &frontDoorProfileDeletePoller{}

type frontDoorProfileDeletePoller struct {
	client *profiles.ProfilesClient
	id     profiles.ProfileId
}

func NewFrontDoorProfileDeletePoller(client *profiles.ProfilesClient, id profiles.ProfileId) pollers.PollerType {
	return &frontDoorProfileDeletePoller{
		client: client,
		id:     id,
	}
}

func (p frontDoorProfileDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	err := p.client.DeleteThenPoll(ctx, p.id)
	if err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "another operation is in progress") || strings.Contains(msg, "operation is in progress") {
			return &pollers.PollResult{
				PollInterval: 30 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}

		return nil, fmt.Errorf("deleting %s: %+v", p.id, err)
	}

	return &pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
