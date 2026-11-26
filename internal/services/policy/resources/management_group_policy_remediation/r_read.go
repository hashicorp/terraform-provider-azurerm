// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_policy_remediation

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupPolicyRemediationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseProviders2RemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("reading Policy Remediation: %+v", err)
	}

	resp, err := client.GetAtManagementGroup(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id.ID(), err)
	}

	d.Set("name", id.RemediationName)
	managementGroupID := parse.NewManagementGroupId(id.ManagementGroupId)
	d.Set("management_group_id", managementGroupID.ID())

	if props := resp.Model.Properties; props != nil {
		locations := make([]any, 0)
		if filters := props.Filters; filters != nil {
			locations = pluginsdk.FlattenSlice(filters.Locations)
		}
		if err := d.Set("location_filters", locations); err != nil {
			return fmt.Errorf("setting `location_filters`: %+v", err)
		}

		d.Set("policy_assignment_id", props.PolicyAssignmentId)
		d.Set("policy_definition_reference_id", props.PolicyDefinitionReferenceId)

		d.Set("resource_count", props.ResourceCount)
		d.Set("parallel_deployments", props.ParallelDeployments)
		if props.FailureThreshold != nil {
			d.Set("failure_percentage", props.FailureThreshold.Percentage)
		}
	}

	return nil
}

func managementGroupPolicyRemediationCancellationRefreshFunc(ctx context.Context,
	client *remediations.RemediationsClient, id remediations.Providers2RemediationId,
) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
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
