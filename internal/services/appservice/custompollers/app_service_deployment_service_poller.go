// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &appServiceDeploymentServicePoller{}

type appServiceDeploymentServicePoller struct {
	statusReq *http.Request
}

func NewAppServiceDeploymentServicePoller(statusReq *http.Request) *appServiceDeploymentServicePoller {
	return &appServiceDeploymentServicePoller{
		statusReq: statusReq,
	}
}

func (p appServiceDeploymentServicePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	reqWithTimeout := p.statusReq.Clone(attemptCtx)

	resp, err := http.DefaultClient.Do(reqWithTimeout)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("[DEBUG] Deployment service isn't available yet, retrying...")
			return &pollers.PollResult{
				Status:       pollers.PollingStatusInProgress,
				PollInterval: 30 * time.Second,
			}, nil
		}
		return nil, fmt.Errorf("client error: %s", err)
	}

	if resp.StatusCode >= 500 {
		resp.Body.Close()
		log.Printf("[DEBUG] Deployment service came back with a %d status code, retrying...", resp.StatusCode)
		return &pollers.PollResult{
			Status:       pollers.PollingStatusInProgress,
			PollInterval: 30 * time.Second,
		}, nil
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("client error: %s", resp.Status)
	}

	return &pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 30 * time.Second,
	}, nil
}
