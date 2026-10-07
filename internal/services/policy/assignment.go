// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
)

func convertEnforcementMode(mode bool) *policyassignments.EnforcementMode {
	m := policyassignments.EnforcementModeDoNotEnforce
	if mode {
		m = policyassignments.EnforcementModeDefault
	}
	return &m
}

func waitForPolicyAssignmentToStabilize(ctx context.Context, client *policyassignments.PolicyAssignmentsClient, id policyassignments.ScopedPolicyAssignmentId, shouldExist bool) error {
	pollerOpts := &custompollers.EventualConsistencyPollerOptions{
		Interval:              time.Second * 5,
		TargetStatusCode:      pointer.To(http.StatusOK),
		RetryErrorStatusCodes: []int{http.StatusNotFound},
	}

	if !shouldExist {
		pollerOpts.TargetStatusCode = pointer.To(http.StatusNotFound)
		pollerOpts.RetryErrorStatusCodes = nil
	}

	poller := custompollers.NewEventualConsistencyPoller(20, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, id)
		return resp.HttpResponse, err
	}, pollerOpts)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling for %s: %+v", id, err)
	}

	return nil
}
