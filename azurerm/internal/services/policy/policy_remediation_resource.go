package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/policyinsights/mgmt/2019-10-01-preview/policyinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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
			_, err := parse.PolicyRemediationID(id)
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
				ValidateFunc:     validate.PolicyScopeID,
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
			},
		},
	}
}

func resourceArmPolicyRemediationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope, err := parse.PolicyScopeID(d.Get("scope").(string))
	if err != nil {
		return fmt.Errorf("creating/updating Policy Remediation %q: %+v", name, err)
	}

	if d.IsNewResource() {
		existing, err := RemediationGetAtScope(ctx, client, name, scope)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Policy Remediation %q (Scope %q): %+v", name, scope.ScopeId(), err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_remediation", *existing.ID)
		}
	}

	parameters := policyinsights.Remediation{
		RemediationProperties: &policyinsights.RemediationProperties{
			Filters: &policyinsights.RemediationFilters{
				Locations: utils.ExpandStringSlice(d.Get("location_filters").([]interface{})),
			},
			PolicyAssignmentID:          utils.String(d.Get("policy_assignment_id").(string)),
			PolicyDefinitionReferenceID: utils.String(d.Get("policy_definition_reference_id").(string)),
		},
	}

	switch scope := scope.(type) {
	case parse.ScopeAtSubscription:
		_, err = client.CreateOrUpdateAtSubscription(ctx, scope.SubscriptionId, name, parameters)
	case parse.ScopeAtResourceGroup:
		_, err = client.CreateOrUpdateAtResourceGroup(ctx, scope.SubscriptionId, scope.ResourceGroup, name, parameters)
	case parse.ScopeAtResource:
		_, err = client.CreateOrUpdateAtResource(ctx, scope.ScopeId(), name, parameters)
	case parse.ScopeAtManagementGroup:
		_, err = client.CreateOrUpdateAtManagementGroup(ctx, scope.ManagementGroupName, name, parameters)
	default:
		return fmt.Errorf("creating/updating Policy Remediation %q: invalid scope type", name)
	}
	if err != nil {
		return fmt.Errorf("creating/updating Policy Remediation %q (Scope %q): %+v", name, scope.ScopeId(), err)
	}

	resp, err := RemediationGetAtScope(ctx, client, name, scope)
	if err != nil {
		return fmt.Errorf("retrieving Policy Remediation %q (Scope %q): %+v", name, scope.ScopeId(), err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Policy Remediation %q (Scope %q)", name, scope.ScopeId())
	}
	d.SetId(*resp.ID)

	return resourceArmPolicyRemediationRead(d, meta)
}

func resourceArmPolicyRemediationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyRemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("reading Policy Remediation: %+v", err)
	}

	resp, err := RemediationGetAtScope(ctx, client, id.Name, id.PolicyScopeId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Policy Remediation %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Policy Remediation %q (Scope %q): %+v", id.Name, id.ScopeId(), err)
	}

	d.Set("name", id.Name)
	d.Set("scope", id.ScopeId())

	if props := resp.RemediationProperties; props != nil {
		locations := []interface{}{}
		if filters := props.Filters; filters != nil {
			locations = utils.FlattenStringSlice(filters.Locations)
		}
		if err := d.Set("location_filters", locations); err != nil {
			return fmt.Errorf("setting `location_filters`: %+v", err)
		}

		d.Set("policy_assignment_id", props.PolicyAssignmentID)
		d.Set("policy_definition_reference_id", props.PolicyDefinitionReferenceID)
	}

	return nil
}

func resourceArmPolicyRemediationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyRemediationID(d.Id())
	if err != nil {
		return err
	}

	switch scope := id.PolicyScopeId.(type) {
	case parse.ScopeAtSubscription:
		_, err = client.DeleteAtSubscription(ctx, scope.SubscriptionId, id.Name)
	case parse.ScopeAtResourceGroup:
		_, err = client.DeleteAtResourceGroup(ctx, scope.SubscriptionId, scope.ResourceGroup, id.Name)
	case parse.ScopeAtResource:
		_, err = client.DeleteAtResource(ctx, scope.ScopeId(), id.Name)
	case parse.ScopeAtManagementGroup:
		_, err = client.DeleteAtManagementGroup(ctx, scope.ManagementGroupName, id.Name)
	default:
		return fmt.Errorf("deleting Policy Remediation %q: invalid scope type", id.Name)
	}
	if err != nil {
		return fmt.Errorf("deleting Policy Remediation %q (Scope %q): %+v", id.Name, id.ScopeId(), err)
	}

	return nil
}

// RemediationGetAtScope is a wrapper of the 4 Get functions on RemediationsClient, combining them into one to simplify code.
func RemediationGetAtScope(ctx context.Context, client *policyinsights.RemediationsClient, name string, scopeId parse.PolicyScopeId) (policyinsights.Remediation, error) {
	switch scopeId := scopeId.(type) {
	case parse.ScopeAtSubscription:
		return client.GetAtSubscription(ctx, scopeId.SubscriptionId, name)
	case parse.ScopeAtResourceGroup:
		return client.GetAtResourceGroup(ctx, scopeId.SubscriptionId, scopeId.ResourceGroup, name)
	case parse.ScopeAtResource:
		return client.GetAtResource(ctx, scopeId.ScopeId(), name)
	case parse.ScopeAtManagementGroup:
		return client.GetAtManagementGroup(ctx, scopeId.ManagementGroupName, name)
	default:
		return policyinsights.Remediation{}, fmt.Errorf("reading Policy Remediation %q: invalid scope type", name)
	}
}
