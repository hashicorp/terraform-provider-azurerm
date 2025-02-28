// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	managmentGroupParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	managmentGroupValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmManagementGroupPolicyRemediation() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceArmManagementGroupPolicyRemediationCreateUpdate,
		Read:   resourceArmManagementGroupPolicyRemediationRead,
		Update: resourceArmManagementGroupPolicyRemediationCreateUpdate,
		Delete: resourceArmManagementGroupPolicyRemediationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ResourcePolicyRemediationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RemediationName,
			},

			"management_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: managmentGroupValidate.ManagementGroupID,
			},

			"policy_assignment_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.PolicyAssignmentID,
			},

			"failure_percentage": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(0, 1.0),
			},

			"parallel_deployments": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntPositive,
			},

			"resource_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntPositive,
			},

			"location_filters": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: location.EnhancedValidate,
				},
			},

			"policy_definition_reference_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}

	return resource
}

func resourceArmManagementGroupPolicyRemediationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementID, err := managmentGroupParse.ManagementGroupID(d.Get("management_group_id").(string))
	if err != nil {
		return err
	}
	id := remediations.NewProviders2RemediationID(managementID.Name, d.Get("name").(string))

	if d.IsNewResource() {
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

	var parameters remediations.Remediation
	props := &remediations.RemediationProperties{
		Filters: &remediations.RemediationFilters{
			Locations: utils.ExpandStringSlice(d.Get("location_filters").([]interface{})),
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

	d.SetId(id.ID())

	return resourceArmManagementGroupPolicyRemediationRead(d, meta)
}

func resourceArmManagementGroupPolicyRemediationRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	managementGroupID := managmentGroupParse.NewManagementGroupId(id.ManagementGroupId)
	d.Set("management_group_id", managementGroupID.ID())

	if props := resp.Model.Properties; props != nil {
		locations := make([]interface{}, 0)
		if filters := props.Filters; filters != nil {
			locations = utils.FlattenStringSlice(filters.Locations)
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

func resourceArmManagementGroupPolicyRemediationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseProviders2RemediationID(d.Id())
	if err != nil {
		return err
	}

	// we have to cancel the remediation first before deleting it when the resource_discovery_mode is set to ReEvaluateCompliance
	// therefore we first retrieve the remediation to see if the resource_discovery_mode is switched to ReEvaluateCompliance
	existing, err := client.GetAtManagementGroup(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := waitForRemediationToDelete(ctx, existing.Model.Properties, id.ID(), d.Timeout(pluginsdk.TimeoutDelete),
		func() error {
			_, err := client.CancelAtManagementGroup(ctx, *id)
			return err
		},
		managementGroupPolicyRemediationCancellationRefreshFunc(ctx, client, *id),
	); err != nil {
		return err
	}

	_, err = client.DeleteAtManagementGroup(ctx, *id)

	return err
}

func managementGroupPolicyRemediationCancellationRefreshFunc(ctx context.Context,
	client *remediations.RemediationsClient, id remediations.Providers2RemediationId,
) pluginsdk.StateRefreshFunc {
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
