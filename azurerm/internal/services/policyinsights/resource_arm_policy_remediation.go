package policyinsights

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/policyinsights/mgmt/2019-10-01/policyinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicyInsightsRemediation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicyInsightsRemediationCreateUpdate,
		Read:   resourceArmPolicyInsightsRemediationRead,
		Update: resourceArmPolicyInsightsRemediationCreateUpdate,
		Delete: resourceArmPolicyInsightsRemediationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				ValidateFunc: validateRemediationName,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"policy_assignment_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePolicyAssignmentID,
			},

			"policy_definition_reference_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePolicyDefinitionID,
			},

			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_updated_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmPolicyInsightsRemediationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PolicyInsights.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	scopeObj, err := ParseScope(scope)
	if err != nil {
		return fmt.Errorf("Error creating Policy Remediation %q: %+v", name, err)
	}

	switch scopeObj.Type {
	case AtSubscription:
		err = remediationCreateUpdateAtSubscription(client, ctx, scopeObj, d)
	case AtManagementGroup:
		err = remediationCreateUpdateAtManagementGroup(client, ctx, scopeObj, d)
	case AtResourceGroup:
		err = remediationCreateUpdateAtResourceGroup(client, ctx, scopeObj, d)
	case AtResource:
		err = remediationCreateUpdateAtResource(client, ctx, scopeObj, d)
	default:
		return fmt.Errorf("Error creating Policy Remediation %q: Cannot recognize scope %q as Subscription ID, Management Group ID, Resource Group ID, or Resource ID", name, scope)
	}

	if err != nil {
		return err
	}

	return resourceArmPolicyInsightsRemediationRead(d, meta)
}

func resourceArmPolicyInsightsRemediationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PolicyInsights.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	remediationId, err := ParseRemediationId(d.Id())
	if err != nil {
		return fmt.Errorf("Error reading Policy Insight Remediation: %+v", err)
	}

	name := remediationId.Name
	scope := remediationId.Scope

	d.Set("name", name)
	d.Set("scope", scope)

	var resp policyinsights.Remediation

	switch remediationId.Type {
	case AtSubscription:
		resp, err = client.GetAtSubscription(ctx, *remediationId.SubscriptionId, name)
	case AtManagementGroup:
		resp, err = client.GetAtManagementGroup(ctx, *remediationId.ManagementGroupId, name)
	case AtResourceGroup:
		resp, err = client.GetAtResourceGroup(ctx, *remediationId.SubscriptionId, *remediationId.ResourceGroup, name)
	case AtResource:
		resp, err = client.GetAtResource(ctx, scope, name)
	default:
		return fmt.Errorf("Error reading Policy Remediation %q: Cannot recognize scope %q as Subscription ID, Management Group ID, Resource Group ID, or Resource ID", name, scope)
	}
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Policy Remediation %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Policy Remediation %q (Scope %q): %+v", name, scope, err)
	}

	if props := resp.RemediationProperties; props != nil {
		if err := d.Set("location_filters", flattenArmRemediationLocationFilters(props.Filters)); err != nil {
			return fmt.Errorf("Error setting `location_filters`: %+v", err)
		}

		d.Set("created_on", (props.CreatedOn).String())
		d.Set("last_updated_on", (props.LastUpdatedOn).String())
		d.Set("policy_assignment_id", props.PolicyAssignmentID)
		d.Set("policy_definition_reference_id", props.PolicyDefinitionReferenceID)
	}

	return nil
}

func resourceArmPolicyInsightsRemediationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PolicyInsights.RemediationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	remediationId, err := ParseRemediationId(d.Id())
	if err != nil {
		return err
	}

	name := remediationId.Name
	scope := remediationId.Scope

	switch remediationId.Type {
	case AtSubscription:
		_, err = client.DeleteAtSubscription(ctx, *remediationId.SubscriptionId, name)
	case AtManagementGroup:
		_, err = client.DeleteAtManagementGroup(ctx, *remediationId.ManagementGroupId, name)
	case AtResourceGroup:
		_, err = client.DeleteAtResourceGroup(ctx, *remediationId.SubscriptionId, *remediationId.ResourceGroup, name)
	case AtResource:
		_, err = client.DeleteAtResource(ctx, scope, name)
	default:
		return fmt.Errorf("Error deleting Policy Insight Remediation %q: Cannot recognize scope %q as Subscription ID, Management Group ID, Resource Group ID, or Resource ID", name, scope)
	}

	if err != nil {
		return err
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

	var result []interface{}
	for _, location := range *input.Locations {
		result = append(result, location)
	}

	return result
}
