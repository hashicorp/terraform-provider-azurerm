// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/appserviceplans"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &servicePlanDeletePoller{}

type servicePlanDeletePoller struct {
	client  *appserviceplans.AppServicePlansClient
	id      commonids.AppServicePlanId
	retries int
}

func NewServicePlanDeletePoller(client *appserviceplans.AppServicePlansClient, id commonids.AppServicePlanId) *servicePlanDeletePoller {
	return &servicePlanDeletePoller{
		client: client,
		id:     id,
	}
}

func (p *servicePlanDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	resp, err := p.client.Delete(ctx, p.id)
	if err == nil || response.WasNotFound(resp.HttpResponse) {
		return &pollers.PollResult{Status: pollers.PollingStatusSucceeded}, nil
	}

	// Deleted apps can remain associated with the plan briefly. Retry this conflict up to 10 times.
	if response.WasConflict(resp.HttpResponse) && strings.Contains(err.Error(), "contain resources") && p.retries < 10 {
		p.retries++
		return &pollers.PollResult{
			Status:       pollers.PollingStatusInProgress,
			PollInterval: 10 * time.Second,
		}, nil
	}

	return nil, pollers.PollingFailedError{Message: err.Error()}
}
