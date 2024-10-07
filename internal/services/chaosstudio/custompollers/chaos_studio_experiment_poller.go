// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &chaosStudioExperimentPoller{}

type chaosStudioExperimentPoller struct {
	client *experiments.ExperimentsClient
	id     experiments.ExperimentId
}

var (
	pollingSuccess = pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 10 * time.Second,
	}
	pollingInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}
)

func NewChaosStudioExperimentPoller(client *experiments.ExperimentsClient, id experiments.ExperimentId) *chaosStudioExperimentPoller {
	return &chaosStudioExperimentPoller{
		client: client,
		id:     id,
	}
}

type operationResult struct {
	Name              *string `json:"name"`
	ProvisioningState string  `json:"provisioningState"`
}

func (p chaosStudioExperimentPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	var respBody []byte
	respBody, err = io.ReadAll(resp.HttpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing response body: %+v", err)
	}
	resp.HttpResponse.Body.Close()
	resp.HttpResponse.Body = io.NopCloser(bytes.NewReader(respBody))

	var op operationResult
	contentType := resp.HttpResponse.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "application/json") {

		if err = json.Unmarshal(respBody, &op); err != nil {
			return nil, fmt.Errorf("unmarshalling response body: %+v", err)
		}
	} else {
		return nil, fmt.Errorf("internal-error: polling support for the Content-Type %q was not implemented: %+v", contentType, err)
	}

	if op.ProvisioningState != string(pollingSuccess.Status) {
		return &pollingInProgress, err
	}
	return &pollingSuccess, nil
}
