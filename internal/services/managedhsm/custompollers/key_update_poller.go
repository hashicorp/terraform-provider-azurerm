// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

var _ pollers.PollerType = &keyUpdatePoller{}

type keyUpdatePoller struct {
	client         *keyvault.BaseClient
	baseURI        string
	keyName        string
	expectedUpdate date.UnixTime
}

func NewKeyUpdatePoller(client *keyvault.BaseClient, baseURI, keyName string, expectedUpdate date.UnixTime) pollers.PollerType {
	return &keyUpdatePoller{
		client:         client,
		baseURI:        baseURI,
		keyName:        keyName,
		expectedUpdate: expectedUpdate,
	}
}

func (p *keyUpdatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.GetKey(ctx, p.baseURI, p.keyName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("key %q was not found after update", p.keyName)
		}
		return nil, fmt.Errorf("polling key %q after update: %+v", p.keyName, err)
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
