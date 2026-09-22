// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceResourceGroupPolicyRemediation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceResourceGroupPolicyRemediationCreateUpdate,
		Read:   resourceResourceGroupPolicyRemediationRead,
		Update: resourceResourceGroupPolicyRemediationCreateUpdate,
		Delete: resourceResourceGroupPolicyRemediationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := remediations.ParseProviderRemediationID(id)
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

			"resource_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.AsGeneratedID(commonids.ParseResourceGroupIDInsensitively),
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
				// The API does not honour the provided casing, instead returning values in all lowercase
				// https://github.com/Azure/azure-rest-api-specs/issues/37823
				DiffSuppressFunc:      suppress.CaseDifference,
				DiffSuppressOnRefresh: true,
			},

			"resource_discovery_mode": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(remediations.ResourceDiscoveryModeExistingNonCompliant),
				ValidateFunc: validation.StringInSlice(remediations.PossibleValuesForResourceDiscoveryMode(), false),
			},
		},
	}
}

func resourceResourceGroupPolicyRemediationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	resourceGroupId, err := commonids.ParseResourceGroupIDInsensitively(d.Get("resource_group_id").(string))
	if err != nil {
		return err
	}

	id := remediations.NewProviderRemediationID(resourceGroupId.SubscriptionId, resourceGroupId.ResourceGroupName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetAtResourceGroup(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
				}
			}
			if existing.Model != nil {
				return tf.ImportAsExistsError("azurerm_resource_group_policy_remediation", id.ID())
			}
		}
	}

	parameters := remediations.Remediation{
		Properties: readRemediationProperties(d),
	}

	if _, err = client.CreateOrUpdateAtResourceGroup(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id.ID(), err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceResourceGroupPolicyRemediationRead(d, meta)
}

func resourceResourceGroupPolicyRemediationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseProviderRemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("reading Policy Remediation: %+v", err)
	}

	resourceGroupId := commonids.NewResourceGroupID(id.SubscriptionId, id.ResourceGroupName)

	resp, err := client.GetAtResourceGroup(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id.ID(), err)
	}

	d.Set("name", id.RemediationName)
	d.Set("resource_group_id", resourceGroupId.ID())

	return setRemediationProperties(d, resp.Model.Properties)
}

func resourceResourceGroupPolicyRemediationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseProviderRemediationID(d.Id())
	if err != nil {
		return err
	}

	// we have to cancel the remediation first before deleting it when the resource_discovery_mode is set to ReEvaluateCompliance
	// therefore we first retrieve the remediation to see if the resource_discovery_mode is switched to ReEvaluateCompliance
	existing, err := client.GetAtResourceGroup(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := waitForRemediationToDelete(
		ctx, existing.Model.Properties, id.ID(), d.Timeout(pluginsdk.TimeoutDelete),
		func() error {
			_, err := client.CancelAtResourceGroup(ctx, *id)
			return err
		},
		resourceGroupPolicyRemediationCancellationRefreshFunc(ctx, client, *id),
	); err != nil {
		return fmt.Errorf("waiting for remediation to delete %s: %+v", id, err)
	}

	_, err = client.DeleteAtResourceGroup(ctx, *id)

	return err
}

func resourceGroupPolicyRemediationCancellationRefreshFunc(ctx context.Context, client *remediations.RemediationsClient, id remediations.ProviderRemediationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetAtResourceGroup(ctx, id)
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
