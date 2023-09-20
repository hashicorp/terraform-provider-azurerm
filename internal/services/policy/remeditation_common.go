package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func deleteRemediation(ctx context.Context, id remediations.ScopedRemediationId, client *remediations.RemediationsClient) error {
	// TODO: refactor this

	// we have to cancel the remediation first before deleting it when the resource_discovery_mode is set to ReEvaluateCompliance
	// therefore we first retrieve the remediation to see if the resource_discovery_mode is switched to ReEvaluateCompliance
	existing, err := client.GetAtResource(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := waitForRemediationToDelete(ctx, existing.Model.Properties, id.ID(), d.Timeout(pluginsdk.TimeoutDelete),
		func() error {
			_, err := client.CancelAtResource(ctx, id)
			return err
		},
		managementGroupPolicyRemediationCancellationRefreshFunc(ctx, client, *id),
	); err != nil {
		return err
	}

	_, err = client.DeleteAtResource(ctx, *id)

	return err
}

// waitForRemediationToDelete waits for the remediation to a status that allow to delete
func waitForRemediationToDelete(ctx context.Context,
	prop *remediations.RemediationProperties,
	id string,
	timeout time.Duration,
	cancelFunc func() error,
	refresh pluginsdk.StateRefreshFunc) error {
	if prop == nil {
		return nil
	}
	if mode := prop.ResourceDiscoveryMode; mode != nil && *mode == remediations.ResourceDiscoveryModeReEvaluateCompliance {
		// Remediation can only be canceld when it is in "Evaluating" or "Accepted" status, otherwise, API might raise error (e.g. canceling a "Completed" remediation returns 400).
		if state := prop.ProvisioningState; state != nil && (*state == "Evaluating" || *state == "Accepted") {
			log.Printf("[DEBUG] cancelling the remediation first before deleting it when `resource_discovery_mode` is set to `ReEvaluateCompliance`")
			if err := cancelFunc(); err != nil {
				return fmt.Errorf("cancelling %s: %+v", id, err)
			}

			log.Printf("[DEBUG] waiting for the %s to be canceled", id)
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"Cancelling"},
				Target: []string{
					"Succeeded", "Canceled", "Failed",
				},
				Refresh:    refresh,
				MinTimeout: 10 * time.Second,
				Timeout:    timeout,
			}

			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be canceled: %+v", id, err)
			}
		}
	}
	return nil
}

func managementGroupPolicyRemediationCancellationRefreshFunc(ctx context.Context,
	client *remediations.RemediationsClient, id remediations.Providers2RemediationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetAtManagementGroup(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("issuing read request for %s: %+v", id.ID(), err)
		}

		if resp.Model.Properties == nil {
			return nil, "", fmt.Errorf("`properties` was nil")
		}
		if resp.Model.Properties.ProvisioningState == nil {
			return nil, "", fmt.Errorf("`properties.ProvisioningState` was nil")
		}
		return resp, *resp.Model.Properties.ProvisioningState, nil
	}
}
