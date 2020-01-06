package policyinsights

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/policyinsights/mgmt/2019-10-01/policyinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPolicyInsightsRemediation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPolicyInsightsRemediationRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location_filters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_updated_on": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"policy_assignment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"policy_definition_reference_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmPolicyInsightsRemediationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PolicyInsights.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	scopeObj, err := ParseScope(scope)
	if err != nil {
		return fmt.Errorf("Error reading Policy Remediation %q (Scope %q): %+v", name, scope, err)
	}

	var resp policyinsights.Remediation

	switch scopeObj.Type {
	case AtSubscription:
		resp, err = client.GetAtSubscription(ctx, *scopeObj.SubscriptionId, name)
	case AtManagementGroup:
		resp, err = client.GetAtManagementGroup(ctx, *scopeObj.ManagementGroupId, name)
	case AtResourceGroup:
		resp, err = client.GetAtResourceGroup(ctx, *scopeObj.SubscriptionId, *scopeObj.ResourceGroup, name)
	case AtResource:
		resp, err = client.GetAtResource(ctx, scopeObj.Scope, name)
	default:
		return fmt.Errorf("Error reading Policy Remediation %q: Cannot recognize scope %q as Subscription ID, Management Group ID, Resource Group ID or Resource ID", name, scope)
	}
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Policy Remediation %q (Scope %q) was not found", name, scope)
		}
		return fmt.Errorf("Error reading Policy Remediation %q (Scope %q): %+v", name, scope, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Policy Remediation %q (Scope %q) ID", name, scope)
	}

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("scope", scope)

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
