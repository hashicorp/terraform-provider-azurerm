package policyinsights

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policyinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policyinsights/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPolicyRemediation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPolicyRemediationRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.RemediationName,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.RemediationScopeID,
			},
		},
	}
}

func dataSourceArmPolicyRemediationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PolicyInsights.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)
	scopeId, err := parse.RemediationScopeID(scope)
	if err != nil {
		return fmt.Errorf("Error reading Policy Remediation %q (Scope %q): %+v", name, scope, err)
	}

	resp, err := RemediationGetAtScope(ctx, client, name, *scopeId)
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

	return nil
}
