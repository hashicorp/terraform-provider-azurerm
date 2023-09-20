// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceManagementGroupRemediation() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceManagementGroupRemediationCreateUpdate,
		Read:   resourceManagementGroupRemediationRead,
		Update: resourceManagementGroupRemediationCreateUpdate,
		Delete: resourceManagementGroupRemediationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ParseManagementGroupRemediationID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ManagementGroupRemediationV0ToV1{},
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

			"management_group_id": commonschema.ResourceIDReferenceRequiredForceNew(commonids.ManagementGroupId{}),

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
				// TODO: add validation to this field
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

	if !features.FourPointOhBeta() {
		resource.Schema["resource_discovery_mode"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(remediations.ResourceDiscoveryModeExistingNonCompliant),
			ValidateFunc: validation.StringInSlice([]string{
				string(remediations.ResourceDiscoveryModeExistingNonCompliant),
				string(remediations.ResourceDiscoveryModeReEvaluateCompliance),
			}, false),
			Deprecated: "`resource_discovery_mode` will be removed in version 4.0 of the AzureRM Provider as evaluating compliance before remediation is only supported at subscription scope and below.",
		}
	}
	return resource
}

func resourceManagementGroupRemediationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementGroupId, err := commonids.ParseManagementGroupID(d.Get("management_group_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewManagementGroupRemediationID(*managementGroupId, d.Get("name").(string)).ToRemediationID()
	if d.IsNewResource() {
		existing, err := client.GetAtResource(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_management_group_policy_remediation", id.ID())
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
	return resourceManagementGroupRemediationRead(d, meta)
}

func resourceManagementGroupRemediationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsed, err := parse.ParseManagementGroupRemediationID(d.Id())
	if err != nil {
		return err
	}
	id := parsed.ToRemediationID()

	resp, err := client.GetAtResource(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RemediationName)
	d.Set("management_group_id", parsed.ManagementGroupId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			locations := make([]interface{}, 0)
			if filters := props.Filters; filters != nil {
				locations = utils.FlattenStringSlice(filters.Locations)
			}
			if err := d.Set("location_filters", locations); err != nil {
				return fmt.Errorf("setting `location_filters`: %+v", err)
			}

			d.Set("policy_assignment_id", props.PolicyAssignmentId)                    // TODO: normalize
			d.Set("policy_definition_reference_id", props.PolicyDefinitionReferenceId) // TODO: normalize
			d.Set("resource_discovery_mode", string(pointer.From(props.ResourceDiscoveryMode)))

			d.Set("resource_count", props.ResourceCount)
			d.Set("parallel_deployments", props.ParallelDeployments)
			failurePercentage := 0.0
			if props.FailureThreshold != nil {
				failurePercentage = pointer.From(props.FailureThreshold.Percentage)
			}
			d.Set("failure_percentage", failurePercentage)
		}
	}

	return nil
}

func resourceManagementGroupRemediationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseManagementGroupRemediationID(d.Id())
	if err != nil {
		return err
	}

	return deleteRemediation(ctx, id.ToRemediationID(), client)
}
