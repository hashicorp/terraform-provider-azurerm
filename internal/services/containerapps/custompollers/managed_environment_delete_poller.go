// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type managedEnvironmentReader interface {
	Get(context.Context, managedenvironments.ManagedEnvironmentId) (managedenvironments.GetOperationResponse, error)
}

type managedEnvironmentDeletePoller struct {
	client managedEnvironmentReader
	id     managedenvironments.ManagedEnvironmentId
}

func NewManagedEnvironmentDeletePoller(client managedEnvironmentReader, id managedenvironments.ManagedEnvironmentId) pollers.Poller {
	return pollers.NewPoller(&managedEnvironmentDeletePoller{client: client, id: id}, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
}

func (p *managedEnvironmentDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	result, err := p.client.Get(ctx, p.id)
	if response.WasNotFound(result.HttpResponse) {
		return &pollers.PollResult{Status: pollers.PollingStatusSucceeded}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("checking deletion of %s: %w", p.id, err)
	}
	if result.HttpResponse == nil {
		return nil, fmt.Errorf("checking deletion of %s: response was nil", p.id)
	}
	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}, nil
}
