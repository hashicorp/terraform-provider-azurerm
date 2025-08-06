// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverazureadonlyauthentications"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &MsSqlServerDeleteServerAzureADOnlyAuthenticationPoller{}

func NewMsSqlServerDeleteServerAzureADOnlyAuthenticationPoller(client *serverazureadonlyauthentications.ServerAzureADOnlyAuthenticationsClient, serverId commonids.SqlServerId) *MsSqlServerDeleteServerAzureADOnlyAuthenticationPoller {
	return &MsSqlServerDeleteServerAzureADOnlyAuthenticationPoller{
		client:   client,
		serverId: serverId,
	}
}

type MsSqlServerDeleteServerAzureADOnlyAuthenticationPoller struct {
	client   *serverazureadonlyauthentications.ServerAzureADOnlyAuthenticationsClient
	serverId commonids.SqlServerId
}

func (p *MsSqlServerDeleteServerAzureADOnlyAuthenticationPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.serverId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.serverId, err)
	}

	if resp.Model == nil {
		return nil, fmt.Errorf("retrieving %s: `model` was nil", p.serverId)
	}

	name := pointer.From(resp.Model.Name)
	azureAdOnlyAuthentication := true
	if props := resp.Model.Properties; props != nil {
		azureAdOnlyAuthentication = props.AzureADOnlyAuthentication
	}

	if !azureAdOnlyAuthentication && strings.EqualFold(name, "Default") {
		return &pollers.PollResult{
			Status:       pollers.PollingStatusSucceeded,
			PollInterval: 5 * time.Second,
		}, nil
	}

	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 5 * time.Second,
	}, nil
}
