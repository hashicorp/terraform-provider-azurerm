// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/verifierworkspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &networkManagerVerifierWorkspacePoller{}

type networkManagerVerifierWorkspacePoller struct {
	client *verifierworkspaces.VerifierWorkspacesClient
	id     verifierworkspaces.VerifierWorkspaceId
}

func NewNetworkManagerVerifierWorkspacePoller(client *verifierworkspaces.VerifierWorkspacesClient, id verifierworkspaces.VerifierWorkspaceId) *networkManagerVerifierWorkspacePoller {
	return &networkManagerVerifierWorkspacePoller{
		client: client,
		id:     id,
	}
}

func (p networkManagerVerifierWorkspacePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollingSuccess, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &pollingInProgress, nil
}
