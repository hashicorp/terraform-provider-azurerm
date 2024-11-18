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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmResourcePolicyRemediation() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceArmResourcePolicyRemediationCreateUpdate,
		Read:   resourceArmResourcePolicyRemediationRead,
		Update: resourceArmResourcePolicyRemediationCreateUpdate,
		Delete: resourceArmResourcePolicyRemediationDelete,

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

			"resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
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

			"resource_discovery_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(remediations.ResourceDiscoveryModeExistingNonCompliant),
				ValidateFunc: validation.StringInSlice([]string{
					string(remediations.ResourceDiscoveryModeExistingNonCompliant),
					string(remediations.ResourceDiscoveryModeReEvaluateCompliance),
				}, false),
			},
		},
	}

	return resource
}

func resourceArmResourcePolicyRemediationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := d.Get("resource_id").(string)

	id := remediations.NewScopedRemediationID(resourceId, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.GetAtResource(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
			}
		}
		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_resource_policy_remediation", id.ID())
		}
	}

	parameters := remediations.Remediation{
		Properties: readRemediationProperties(d),
	}

	if _, err := client.CreateOrUpdateAtResource(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id.ID(), err)
	}

	d.SetId(id.ID())

	return resourceArmResourcePolicyRemediationRead(d, meta)
}

func resourceArmResourcePolicyRemediationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseScopedRemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Policy Scoped Remediation ID: %+v", err)
	}
	resp, err := client.GetAtResource(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id.ID(), err)
	}

	d.Set("name", id.RemediationName)
	d.Set("resource_id", id.ResourceId)

	return setRemediationProperties(d, resp.Model.Properties)
}

func resourceArmResourcePolicyRemediationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseScopedRemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Policy Scoped Remediation ID: %+v", err)
	}

	// we have to cancel the remediation first before deleting it when the resource_discovery_mode is set to ReEvaluateCompliance
	// therefore we first retrieve the remediation to see if the resource_discovery_mode is switched to ReEvaluateCompliance
	existing, err := client.GetAtResource(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := waitForRemediationToDelete(ctx,
		existing.Model.Properties,
		id.ID(),
		d.Timeout(pluginsdk.TimeoutDelete),
		func() error {
			_, err := client.CancelAtResource(ctx, *id)
			return err
		},
		resourcePolicyRemediationCancellationRefreshFunc(ctx, client, *id),
	); err != nil {
		return err
	}

	_, err = client.DeleteAtResource(ctx, *id)

	return err
}

func resourcePolicyRemediationCancellationRefreshFunc(ctx context.Context, client *remediations.RemediationsClient, id remediations.ScopedRemediationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetAtResource(ctx, id)
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

// waitForRemediationToDelete waits for the remediation to a status that allow to delete
func waitForRemediationToDelete(ctx context.Context,
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

// readRemediationProperties sets the properties of the remediation, useful when add new properties to the model
func readRemediationProperties(d *pluginsdk.ResourceData) (prop *remediations.RemediationProperties) {
	prop = &remediations.RemediationProperties{
		Filters: &remediations.RemediationFilters{
			Locations: utils.ExpandStringSlice(d.Get("location_filters").([]interface{})),
		},
		PolicyAssignmentId:          utils.String(d.Get("policy_assignment_id").(string)),
		PolicyDefinitionReferenceId: utils.String(d.Get("policy_definition_reference_id").(string)),
	}
	mode := remediations.ResourceDiscoveryMode(d.Get("resource_discovery_mode").(string))
	prop.ResourceDiscoveryMode = &mode
	if v := d.Get("failure_percentage").(float64); v != 0 {
		prop.FailureThreshold = &remediations.RemediationPropertiesFailureThreshold{
			Percentage: utils.Float(v),
		}
	}
	if v := d.Get("parallel_deployments").(int); v != 0 {
		prop.ParallelDeployments = utils.Int64(int64(v))
	}
	if v := d.Get("resource_count").(int); v != 0 {
		prop.ResourceCount = utils.Int64(int64(v))
	}
	return
}

// setRemediationProperties sets the properties of the remediation, useful when add new properties to the model
func setRemediationProperties(d *pluginsdk.ResourceData, prop *remediations.RemediationProperties) error {
	if prop == nil {
		return nil
	}
	locations := []interface{}{}
	if filters := prop.Filters; filters != nil {
		locations = utils.FlattenStringSlice(filters.Locations)
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
