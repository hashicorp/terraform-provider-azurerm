// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &frontDoorProfileCreatePoller{}

type frontDoorProfileCreatePoller struct {
	client          *profiles.ProfilesClient
	id              profiles.ProfileId
	payload         profiles.Profile
	operationIssued bool
}

func NewFrontDoorProfileCreatePoller(client *profiles.ProfilesClient, id profiles.ProfileId, payload profiles.Profile) pollers.PollerType {
	return &frontDoorProfileCreatePoller{
		client:  client,
		id:      id,
		payload: payload,
	}
}

func (p *frontDoorProfileCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if !p.operationIssued {
		err := p.client.CreateThenPoll(ctx, p.id, p.payload)
		if err != nil {
			if frontDoorProfileOperationInProgress(err) {
				return &pollers.PollResult{
					PollInterval: 30 * time.Second,
					Status:       pollers.PollingStatusInProgress,
				}, nil
			}

			return nil, fmt.Errorf("creating %s: %+v", p.id, err)
		}

		p.operationIssued = true
	}

	ready, err := frontDoorProfileSettled(ctx, p.client, p.id)
	if err != nil {
		return nil, err
	}
	if !ready {
		return &pollers.PollResult{
			PollInterval: 30 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	return &pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}

func frontDoorProfileSettled(ctx context.Context, client *profiles.ProfilesClient, id profiles.ProfileId) (bool, error) {
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return false, nil
		}

		return false, fmt.Errorf("retrieving %s while waiting for profile state to settle: %+v", id, err)
	}

	if resp.Model == nil {
		return false, nil
	}

	return frontDoorProfileStatusesSettled(resp.Model.Properties)
}

func frontDoorProfileStatusesSettled(props *profiles.ProfileProperties) (bool, error) {
	if props == nil || props.ProvisioningState == nil {
		return false, nil
	}

	if *props.ProvisioningState == profiles.ProfileProvisioningStateFailed {
		return false, fmt.Errorf("profile entered failed state with `provisioningState` `%s`", *props.ProvisioningState)
	}

	return *props.ProvisioningState == profiles.ProfileProvisioningStateSucceeded, nil
}

func frontDoorProfileOperationInProgress(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "another operation is in progress") || strings.Contains(msg, "operation is in progress")
}
