// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &CloudHsmClusterStatePoller{}

// CloudHsmClusterStatePoller polls until the Cloud HSM cluster has finished provisioning
// https://github.com/Azure/azure-rest-api-specs/issues/36393
type CloudHsmClusterStatePoller struct {
	client    *cloudhsmclusters.CloudHsmClustersClient
	clusterId cloudhsmclusters.CloudHsmClusterId
}

func NewCloudHsmClusterStatePoller(hsmClient *cloudhsmclusters.CloudHsmClustersClient, clusterId cloudhsmclusters.CloudHsmClusterId) *CloudHsmClusterStatePoller {
	return &CloudHsmClusterStatePoller{
		client:    hsmClient,
		clusterId: clusterId,
	}
}

func (p *CloudHsmClusterStatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.clusterId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.clusterId, err)
	}

	statusMessage := ""
	if model := resp.Model; model != nil && model.Properties != nil {
		statusMessage = strings.ToLower(strings.TrimSpace(pointer.From(model.Properties.StatusMessage)))
	}

	// Check if provisioning is successful
	// HSM cluster update successful or HSM cluster provisioning successful
	if statusMessage == "" || strings.Contains(statusMessage, "success") {
		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			Status: pollers.PollingStatusSucceeded,
		}, nil
	}

	// Check for failure states
	if strings.Contains(statusMessage, "failed") || strings.Contains(statusMessage, "error") {
		return nil, pollers.PollingFailedError{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			Message: fmt.Sprintf("Cloud HSM cluster provisioning failed with status: %s", pointer.From(resp.Model.Properties.StatusMessage)),
		}
	}

	// If we get here, provisioning is still in progress
	return &pollers.PollResult{
		HttpResponse: &client.Response{
			Response: resp.HttpResponse,
		},
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
