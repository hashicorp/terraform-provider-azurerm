// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_policy_remediation

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupPolicyRemediationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementID, err := parse.ManagementGroupID(d.Get("management_group_id").(string))
	if err != nil {
		return err
	}
	id := remediations.NewProviders2RemediationID(managementID.Name, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetAtManagementGroup(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
				}
			}
			if existing.Model != nil {
				return tf.ImportAsExistsError("azurerm_management_group_policy_remediation", id.ID())
			}
		}
	}

	var parameters remediations.Remediation
	props := &remediations.RemediationProperties{
		Filters: &remediations.RemediationFilters{
			Locations: pluginsdk.ExpandStringSlice(d.Get("location_filters").([]any)),
		},
		PolicyAssignmentId:          pointer.To(d.Get("policy_assignment_id").(string)),
		PolicyDefinitionReferenceId: pointer.To(d.Get("policy_definition_reference_id").(string)),
	}

	if v := d.Get("failure_percentage").(float64); v != 0 {
		props.FailureThreshold = &remediations.RemediationPropertiesFailureThreshold{
			Percentage: pointer.To(v),
		}
	}
	if v := d.Get("parallel_deployments").(int); v != 0 {
		props.ParallelDeployments = pointer.To(int64(v))
	}
	if v := d.Get("resource_count").(int); v != 0 {
		props.ResourceCount = pointer.To(int64(v))
	}

	parameters = remediations.Remediation{
		Properties: props,
	}

	if _, err := client.CreateOrUpdateAtManagementGroup(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id.ID(), err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceManagementGroupPolicyRemediationRead(d, meta)
}
