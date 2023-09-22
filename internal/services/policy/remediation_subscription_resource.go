// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubscriptionPolicyRemediation() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceSubscriptionPolicyRemediationCreateUpdate,
		Read:   resourceSubscriptionPolicyRemediationRead,
		Update: resourceSubscriptionPolicyRemediationCreateUpdate,
		Delete: resourceSubscriptionPolicyRemediationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ParseSubscriptionRemediationID(id)
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

			"subscription_id": commonschema.ResourceIDReferenceRequiredForceNew(commonids.SubscriptionId{}),

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

	if !features.FourPointOhBeta() {
		resource.Schema["policy_definition_id"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validate.PolicyDefinitionID,
			Deprecated:       "`policy_definition_id` will be removed in version 4.0 of the AzureRM Provider in favour of `policy_definition_reference_id`.",
		}
	}
	return resource
}

func resourceSubscriptionPolicyRemediationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId, err := commonids.ParseSubscriptionID(d.Get("subscription_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSubscriptionRemediationID(*subscriptionId, d.Get("name").(string)).ToRemediationID()
	if d.IsNewResource() {
		existing, err := client.GetAtResource(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
			}
		}
		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_subscription_policy_remediation", id.ID())
		}
	}

	payload := remediations.Remediation{
		Properties: &remediations.RemediationProperties{
			Filters: &remediations.RemediationFilters{
				Locations: utils.ExpandStringSlice(d.Get("location_filters").([]interface{})),
			},
			PolicyAssignmentId:          pointer.To(d.Get("policy_assignment_id").(string)),
			PolicyDefinitionReferenceId: pointer.To(d.Get("policy_definition_reference_id").(string)),
			ResourceDiscoveryMode:       pointer.To(remediations.ResourceDiscoveryMode(d.Get("resource_discovery_mode").(string))),
		},
	}
	if v := d.Get("failure_percentage").(float64); v != 0 {
		payload.Properties.FailureThreshold = &remediations.RemediationPropertiesFailureThreshold{
			Percentage: utils.Float(v),
		}
	}
	if v := d.Get("parallel_deployments").(int); v != 0 {
		payload.Properties.ParallelDeployments = utils.Int64(int64(v))
	}
	if v := d.Get("resource_count").(int); v != 0 {
		payload.Properties.ResourceCount = utils.Int64(int64(v))
	}

	if _, err := client.CreateOrUpdateAtResource(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSubscriptionPolicyRemediationRead(d, meta)
}

func resourceSubscriptionPolicyRemediationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseRemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("reading Policy Remediation: %+v", err)
	}

	subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	resp, err := client.GetAtSubscription(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id.ID(), err)
	}

	d.Set("name", id.RemediationName)
	d.Set("subscription_id", subscriptionId.ID())

	return setRemediationProperties(d, resp.Model.Properties)
}

func resourceSubscriptionPolicyRemediationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseRemediationID(d.Id())
	if err != nil {
		return err
	}

	// we have to cancel the remediation first before deleting it when the resource_discovery_mode is set to ReEvaluateCompliance
	// therefore we first retrieve the remediation to see if the resource_discovery_mode is switched to ReEvaluateCompliance
	existing, err := client.GetAtSubscription(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := waitForRemediationToDelete(ctx, existing.Model.Properties, id.ID(), d.Timeout(pluginsdk.TimeoutDelete),
		func() error {
			_, err := client.CancelAtSubscription(ctx, *id)
			return err
		},
		subscriptionPolicyRemediationCancellationRefreshFunc(ctx, client, *id),
	); err != nil {
		return err
	}

	_, err = client.DeleteAtSubscription(ctx, *id)

	return err
}

func subscriptionPolicyRemediationCancellationRefreshFunc(ctx context.Context, client *remediations.RemediationsClient, id remediations.RemediationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetAtSubscription(ctx, id)
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
