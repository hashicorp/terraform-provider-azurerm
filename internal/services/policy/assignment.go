// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	assignments "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func convertEnforcementMode(mode bool) *assignments.EnforcementMode {
	m := assignments.EnforcementModeDoNotEnforce
	if mode {
		m = assignments.EnforcementModeDefault
	}
	return &m
}

func waitForPolicyAssignmentToStabilize(ctx context.Context, client *assignments.PolicyAssignmentsClient, id assignments.ScopedPolicyAssignmentId, shouldExist bool) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context was missing a deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"404"},
		Target:  []string{"200"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return resp, strconv.Itoa(resp.HttpResponse.StatusCode), nil
				}

				return nil, strconv.Itoa(resp.HttpResponse.StatusCode), fmt.Errorf("polling for %s: %+v", id, err)
			}

			return resp, strconv.Itoa(resp.HttpResponse.StatusCode), nil
		},
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 20,
		PollInterval:              5 * time.Second,
		Timeout:                   time.Until(deadline),
	}
	if !shouldExist {
		stateConf.Pending = []string{"200"}
		stateConf.Target = []string{"404"}
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}

	return nil
}
