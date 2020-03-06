package policy

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/policyinsights/mgmt/2019-10-01/policyinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicyRemediation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicyRemediationCreateUpdate,
		Read:   resourceArmPolicyRemediationRead,
		Update: resourceArmPolicyRemediationCreateUpdate,
		Delete: resourceArmPolicyRemediationDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.RemediationID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RemediationName,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.RemediationScopeID,
			},

			"policy_assignment_id": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
				// TODO: use the validation function in azurerm_policy_assignment when implemented
				ValidateFunc: validate.PolicyAssignmentID,
			},

			"location_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"policy_definition_reference_id": {
				Type:     schema.TypeString,
				Optional: true,
				// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
				// TODO: use the validation function in azurerm_policy_definition when implemented
				ValidateFunc: validate.PolicyDefinitionID,
			},
		},
	}
}

func resourceArmPolicyRemediationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PolicyInsights.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope, err := parse.RemediationScopeID(d.Get("scope").(string))
	if err != nil {
		return fmt.Errorf("Error creating Policy Remediation %q: %+v", name, err)
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := RemediationGetAtScope(ctx, client, name, *scope)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Policy Remediation %q (Scope %q): %+v", name, scope.ScopeId, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_remediation", *existing.ID)
		}
	}

	filters := d.Get("location_filters").([]interface{})
	policyAssignmentID := d.Get("policy_assignment_id").(string)
	policyDefinitionReferenceID := d.Get("policy_definition_reference_id").(string)

	parameters := policyinsights.Remediation{
		RemediationProperties: &policyinsights.RemediationProperties{
			Filters:                     expandArmRemediationLocationFilters(filters),
			PolicyAssignmentID:          utils.String(policyAssignmentID),
			PolicyDefinitionReferenceID: utils.String(policyDefinitionReferenceID),
		},
	}

	if _, err := RemediationCreateUpdateAtScope(ctx, client, name, *scope, parameters); err != nil {
		return fmt.Errorf("Error creating Policy Remediation %q (Scope %q): %+v", name, scope.ScopeId, err)
	}

	resp, err := RemediationGetAtScope(ctx, client, name, *scope)
	if err != nil {
		return fmt.Errorf("Error retrieving Policy Remediation %q (Scope %q): %+v", name, scope.ScopeId, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Policy Remediation %q (Scope %q) ID", name, scope.ScopeId)
	}
	d.SetId(*resp.ID)

	return resourceArmPolicyRemediationRead(d, meta)
}

func resourceArmPolicyRemediationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PolicyInsights.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("Error reading Policy Remediation: %+v", err)
	}

	resp, err := RemediationGetAtScope(ctx, client, id.Name, id.RemediationScopeId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Policy Remediation %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Policy Remediation %q (Scope %q): %+v", id.Name, id.ScopeId, err)
	}

	d.Set("name", id.Name)
	d.Set("scope", id.ScopeId)

	if props := resp.RemediationProperties; props != nil {
		if err := d.Set("location_filters", flattenArmRemediationLocationFilters(props.Filters)); err != nil {
			return fmt.Errorf("Error setting `location_filters`: %+v", err)
		}

		d.Set("policy_assignment_id", props.PolicyAssignmentID)
		d.Set("policy_definition_reference_id", props.PolicyDefinitionReferenceID)
	}

	return nil
}

func resourceArmPolicyRemediationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PolicyInsights.RemediationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RemediationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := RemediationDeleteAtScope(ctx, client, id.Name, id.RemediationScopeId); err != nil {
		return fmt.Errorf("Error deleting Policy Remediation %q (Scope %q): %+v", id.Name, id.ScopeId, err)
	}

	return nil
}

func expandArmRemediationLocationFilters(input []interface{}) *policyinsights.RemediationFilters {
	if len(input) == 0 {
		return nil
	}

	result := policyinsights.RemediationFilters{
		Locations: utils.ExpandStringSlice(input),
	}

	return &result
}

func flattenArmRemediationLocationFilters(input *policyinsights.RemediationFilters) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return utils.FlattenStringSlice(input.Locations)
}
