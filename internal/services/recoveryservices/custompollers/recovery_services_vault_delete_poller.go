// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &recoveryServicesVaultDeletePoller{}

type recoveryServicesVaultDeletePoller struct {
	client               *vaults.VaultsClient
	id                   vaults.VaultId
	successCount         int
	requiredSuccessCount int
	lastHttpResponse     *http.Response
}

var (
	pollingDeleteSuccess = pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 10 * time.Second,
	}
	pollingDeleteInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}
)

func NewRecoveryServicesVaultDeletePoller(vaultClient *vaults.VaultsClient, id vaults.VaultId) *recoveryServicesVaultDeletePoller {
	return &recoveryServicesVaultDeletePoller{
		client:               vaultClient,
		id:                   id,
		successCount:         0,
		requiredSuccessCount: 5, // Require 5 consecutive 404s to consider delete successful
	}
}

func (p *recoveryServicesVaultDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)

	if err != nil {
		if resp.HttpResponse != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
			p.successCount++
			p.lastHttpResponse = resp.HttpResponse

			if p.successCount >= p.requiredSuccessCount {
				return &pollingDeleteSuccess, nil
			}
			return &pollingDeleteInProgress, nil
		}
		p.successCount = 0
		return &pollingDeleteInProgress, err
	}

	p.successCount = 0
	if resp.HttpResponse != nil {
		p.lastHttpResponse = resp.HttpResponse
	}
	return &pollingDeleteInProgress, nil
}
