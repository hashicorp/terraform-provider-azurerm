// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
)

var _ pollers.PollerType = &logicAppRuntimeTestPoller{}

type logicAppRuntimeTestPoller struct {
	client *resourcemanager.Client
	id     commonids.AppServiceId
}

func NewLogicAppRuntimeTestPoller(client *resourcemanager.Client, id commonids.AppServiceId) *logicAppRuntimeTestPoller {
	return &logicAppRuntimeTestPoller{
		client: client,
		id:     id,
	}
}

func (p logicAppRuntimeTestPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	requestCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	// Handle HTTP 408 in the poller without changing retries for other WebApps requests.
	hostClient := *p.client.Client
	hostClient.DisableRetries = true

	// Portal's Runtime version field uses this endpoint, not the site's provisioning state.
	req, err := hostClient.NewRequest(requestCtx, client.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{http.StatusOK},
		HttpMethod:          http.MethodGet,
		Path:                p.id.ID() + "/host/default/properties/status?api-version=2025-05-01",
	})
	if err != nil {
		return nil, fmt.Errorf("building runtime host status request: %+v", err)
	}

	resp, err := req.Execute(requestCtx)
	if resp != nil && resp.Response != nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("retrieving runtime host status: %+v", err)
		}
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	result := &pollers.PollResult{
		HttpResponse: resp,
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
	if err != nil {
		log.Printf("[DEBUG] Runtime host of %s is not available yet: %+v", p.id, err)
		return result, nil
	}

	var status struct {
		Properties struct {
			State   string `json:"state"`
			Version string `json:"version"`
		} `json:"properties"`
	}
	if err := resp.Unmarshal(&status); err != nil {
		return nil, fmt.Errorf("decoding runtime host status: %+v", err)
	}
	if status.Properties.State == "Running" && status.Properties.Version != "" {
		result.Status = pollers.PollingStatusSucceeded
	} else {
		log.Printf("[DEBUG] Runtime host of %s is not ready: state=%q version=%q", p.id, status.Properties.State, status.Properties.Version)
	}
	return result, nil
}
