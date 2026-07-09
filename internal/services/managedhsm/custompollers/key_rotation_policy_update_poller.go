// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

var _ pollers.PollerType = &keyRotationPolicyUpdatePoller{}

type keyRotationPolicyUpdatePoller struct {
	client         *keyvault.BaseClient
	baseURI        string
	keyName        string
	expectedUpdate date.UnixTime
}

func NewKeyRotationPolicyUpdatePoller(client *keyvault.BaseClient, baseURI, keyName string, expectedUpdate date.UnixTime) pollers.PollerType {
	return &keyRotationPolicyUpdatePoller{
		client:         client,
		baseURI:        baseURI,
		keyName:        keyName,
		expectedUpdate: expectedUpdate,
	}
}

func (p *keyRotationPolicyUpdatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.GetKeyRotationPolicy(ctx, p.baseURI, p.keyName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return nil, fmt.Errorf("key rotation policy for %q was not found after update", p.keyName)
		}
		return nil, fmt.Errorf("polling key rotation policy for %q after update: %+v", p.keyName, err)
	}

	if resp.Attributes == nil || resp.Attributes.Updated == nil {
		return &pollers.PollResult{
			PollInterval: 5 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	if time.Time(*resp.Attributes.Updated).Before(time.Time(p.expectedUpdate)) {
		return &pollers.PollResult{
			PollInterval: 5 * time.Second,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	return &pollers.PollResult{
		PollInterval: 5 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
