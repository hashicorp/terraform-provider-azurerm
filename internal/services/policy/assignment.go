package policy

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func convertEnforcementMode(mode bool) policy.EnforcementMode {
	if mode {
		return policy.EnforcementModeDefault
	} else {
		return policy.EnforcementModeDoNotEnforce
	}
}

func waitForPolicyAssignmentToStabilize(ctx context.Context, client *policy.AssignmentsClient, id parse.PolicyAssignmentId, shouldExist bool) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context was missing a deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"404"},
		Target:  []string{"200"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id.Scope, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return resp, strconv.Itoa(resp.StatusCode), nil
				}

				return nil, strconv.Itoa(resp.StatusCode), fmt.Errorf("polling for %s: %+v", id, err)
			}

			return resp, strconv.Itoa(resp.StatusCode), nil
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
