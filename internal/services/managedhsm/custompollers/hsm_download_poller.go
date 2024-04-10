// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	dataplane "github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

var _ pollers.PollerType = &hsmDownloadPoller{}

func NewHSMDownloadPoller(client *dataplane.HSMSecurityDomainClient, baseUrl string) pollers.PollerType {
	return &hsmDownloadPoller{
		client:  client,
		baseUrl: baseUrl,
	}
}

type hsmDownloadPoller struct {
	client  *dataplane.HSMSecurityDomainClient
	baseUrl string
}

func (p *hsmDownloadPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	res, err := p.client.DownloadPending(ctx, p.baseUrl)
	if err != nil {
		return nil, fmt.Errorf("waiting for Security Domain to download failed within %s: %+v", p.baseUrl, err)
	}

	if res.Status == dataplane.OperationStatusSuccess {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusSucceeded,
			PollInterval: 10 * time.Second,
		}, nil
	}

	// Processing
	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}, nil
}
