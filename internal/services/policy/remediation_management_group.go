package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/policyinsights"
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
	return &pluginsdk.Resource{
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

			"location_filters": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: location.EnhancedValidate,
				},
			},

			"policy_definition_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.PolicyDefinitionID,
			},

			"resource_discovery_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(policyinsights.ResourceDiscoveryModeExistingNonCompliant),
				ValidateFunc: validation.StringInSlice([]string{
					string(policyinsights.ResourceDiscoveryModeExistingNonCompliant),
					string(policyinsights.ResourceDiscoveryModeReEvaluateCompliance),
				}, false),
			},
		},
	}
}

func resourceArmManagementGroupPolicyRemediationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicyInsightsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementID, err := managmentGroupParse.ManagementGroupID(d.Get("management_group_id").(string))
	if err != nil {
		return err
	}
	id := policyinsights.NewProviders2RemediationID(managementID.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.RemediationsGetAtManagementGroup(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
			}
		}
		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_management_group_policy_remediation", *existing.Model.Id)
		}
	}

	parameters := policyinsights.Remediation{
		Properties: &policyinsights.RemediationProperties{
			Filters: &policyinsights.RemediationFilters{
				Locations: utils.ExpandStringSlice(d.Get("location_filters").([]interface{})),
			},
			PolicyAssignmentId:          utils.String(d.Get("policy_assignment_id").(string)),
			PolicyDefinitionReferenceId: utils.String(d.Get("policy_definition_id").(string)),
		},
	}
	mode := policyinsights.ResourceDiscoveryMode(d.Get("resource_discovery_mode").(string))
	parameters.Properties.ResourceDiscoveryMode = &mode

	if _, err := client.RemediationsCreateOrUpdateAtManagementGroup(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id.ID(), err)
	}

	d.SetId(id.ID())

	return resourceArmManagementGroupPolicyRemediationRead(d, meta)
}

func resourceArmManagementGroupPolicyRemediationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicyInsightsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyinsights.ParseProviders2RemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("reading Policy Remediation: %+v", err)
	}

	resp, err := client.RemediationsGetAtManagementGroup(ctx, *id)
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
		locations := []interface{}{}
		if filters := props.Filters; filters != nil {
			locations = utils.FlattenStringSlice(filters.Locations)
		}
		if err := d.Set("location_filters", locations); err != nil {
			return fmt.Errorf("setting `location_filters`: %+v", err)
		}

		d.Set("policy_assignment_id", props.PolicyAssignmentId)
		d.Set("policy_definition_id", props.PolicyDefinitionReferenceId)
		d.Set("resource_discovery_mode", props.ResourceDiscoveryMode)
	}

	return nil
}

func resourceArmManagementGroupPolicyRemediationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicyInsightsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyinsights.ParseProviders2RemediationID(d.Id())
	if err != nil {
		return err
	}

	// we have to cancel the remediation first before deleting it when the resource_discovery_mode is set to ReEvaluateCompliance
	// therefore we first retrieve the remediation to see if the resource_discovery_mode is switched to ReEvaluateCompliance
	existing, err := client.RemediationsGetAtManagementGroup(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := waitRemediationToDelete(ctx, existing.Model.Properties, id.ID(), d.Timeout(pluginsdk.TimeoutDelete),
		func() error {
			_, err := client.RemediationsCancelAtManagementGroup(ctx, *id)
			return err
		},
		managementGroupPolicyRemediationCancellationRefreshFunc(ctx, client, *id),
	); err != nil {
		return err
	}

	_, err = client.RemediationsDeleteAtManagementGroup(ctx, *id)

	return err
}

func managementGroupPolicyRemediationCancellationRefreshFunc(ctx context.Context,
	client *policyinsights.PolicyInsightsClient, id policyinsights.Providers2RemediationId) pluginsdk.StateRefreshFunc {

	return func() (interface{}, string, error) {
		resp, err := client.RemediationsGetAtManagementGroup(ctx, id)
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
