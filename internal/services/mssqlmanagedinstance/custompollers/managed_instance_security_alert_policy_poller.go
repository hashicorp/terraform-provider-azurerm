// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedserversecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &ManagedInstanceSecurityAlertPolicyPoller{}

const managedInstanceSecurityAlertPolicyPollInterval = 30 * time.Second

type ManagedInstanceSecurityAlertPolicyPoller struct {
	client  *managedserversecurityalertpolicies.ManagedServerSecurityAlertPoliciesClient
	id      commonids.SqlManagedInstanceId
	payload managedserversecurityalertpolicies.ManagedServerSecurityAlertPolicy
}

func NewManagedInstanceSecurityAlertPolicyPoller(client *managedserversecurityalertpolicies.ManagedServerSecurityAlertPoliciesClient, id commonids.SqlManagedInstanceId, payload managedserversecurityalertpolicies.ManagedServerSecurityAlertPolicy) *ManagedInstanceSecurityAlertPolicyPoller {
	return &ManagedInstanceSecurityAlertPolicyPoller{
		client:  client,
		id:      id,
		payload: payload,
	}
}

func (p *ManagedInstanceSecurityAlertPolicyPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.CreateOrUpdate(ctx, p.id, p.payload)
	if err != nil {
		// Provisioning a Managed Instance leaves a `Set server security alert policy` operation running
		// server-side, and the API rejects concurrent writes with `ServerSecurityAlertPolicyInProgress`
		// until that completes, so a 409 here means we should send the request again.
		if response.WasConflict(resp.HttpResponse) {
			return &pollers.PollResult{
				PollInterval: managedInstanceSecurityAlertPolicyPollInterval,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}

		return nil, fmt.Errorf("performing CreateOrUpdate: %+v", err)
	}

	if err := resp.Poller.PollUntilDone(ctx); err != nil {
		return nil, fmt.Errorf("polling after CreateOrUpdate: %+v", err)
	}

	return &pollers.PollResult{
		PollInterval: managedInstanceSecurityAlertPolicyPollInterval,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
