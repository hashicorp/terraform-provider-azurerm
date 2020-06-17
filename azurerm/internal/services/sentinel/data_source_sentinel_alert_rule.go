package sentinel

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSentinelAlertRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSentinelAlertRuleRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
			},
		},
	}
}

func dataSourceArmSentinelAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, workspaceID.ResourceGroup, "Microsoft.OperationalInsights", workspaceID.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Sentinel Alert Rule %q (Resource Group %q / Workspace: %q) was not found", name, workspaceID.ResourceGroup, workspaceID.Name)
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule %q (Resource Group %q / Workspace: %q): %+v", name, workspaceID.ResourceGroup, workspaceID.Name, err)
	}

	id := alertRuleID(resp.Value)
	if id == nil || *id == "" {
		return fmt.Errorf("nil or empty ID of Sentinel Alert Rule %q (Resource Group %q / Workspace: %q)", name, workspaceID.ResourceGroup, workspaceID.Name)
	}
	d.SetId(*id)

	return nil
}
