// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// WaitForRemediationToDelete waits for the remediation to a status that allow to delete
func WaitForRemediationToDelete(ctx context.Context,
	prop *remediations.RemediationProperties,
	id string,
	timeout time.Duration,
	cancelFunc func() error,
	refresh pluginsdk.StateRefreshFunc,
) error {
	if prop == nil {
		return nil
	}
	if mode := prop.ResourceDiscoveryMode; mode != nil && *mode == remediations.ResourceDiscoveryModeReEvaluateCompliance {
		// Remediation can only be canceled when it is in "Evaluating" or "Accepted" status, otherwise, API might raise error (e.g. canceling a "Completed" remediation returns 400).
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

// ReadRemediationProperties sets the properties of the remediation, useful when add new properties to the model
func ReadRemediationProperties(d *pluginsdk.ResourceData) (prop *remediations.RemediationProperties) {
	prop = &remediations.RemediationProperties{
		Filters: &remediations.RemediationFilters{
			Locations: pluginsdk.ExpandStringSlice(d.Get("location_filters").([]any)),
		},
		PolicyAssignmentId:          pointer.To(d.Get("policy_assignment_id").(string)),
		PolicyDefinitionReferenceId: pointer.To(d.Get("policy_definition_reference_id").(string)),
	}
	mode := remediations.ResourceDiscoveryMode(d.Get("resource_discovery_mode").(string))
	prop.ResourceDiscoveryMode = &mode
	if v := d.Get("failure_percentage").(float64); v != 0 {
		prop.FailureThreshold = &remediations.RemediationPropertiesFailureThreshold{
			Percentage: pointer.To(v),
		}
	}
	if v := d.Get("parallel_deployments").(int); v != 0 {
		prop.ParallelDeployments = pointer.To(int64(v))
	}
	if v := d.Get("resource_count").(int); v != 0 {
		prop.ResourceCount = pointer.To(int64(v))
	}
	return
}

// SetRemediationProperties sets the properties of the remediation, useful when add new properties to the model
func SetRemediationProperties(d *pluginsdk.ResourceData, prop *remediations.RemediationProperties) error {
	if prop == nil {
		return nil
	}
	locations := []any{}
	if filters := prop.Filters; filters != nil {
		locations = pluginsdk.FlattenSlice(filters.Locations)
	}
	if err := d.Set("location_filters", locations); err != nil {
		return fmt.Errorf("setting `location_filters`: %+v", err)
	}

	d.Set("policy_assignment_id", prop.PolicyAssignmentId)
	d.Set("policy_definition_reference_id", prop.PolicyDefinitionReferenceId)
	d.Set("resource_discovery_mode", pointer.From((*string)(prop.ResourceDiscoveryMode)))

	d.Set("resource_count", prop.ResourceCount)
	d.Set("parallel_deployments", prop.ParallelDeployments)
	if prop.FailureThreshold != nil {
		d.Set("failure_percentage", prop.FailureThreshold.Percentage)
	}
	return nil
}
