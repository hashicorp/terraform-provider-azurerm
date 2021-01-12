package sentinel

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSentinelAlertRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSentinelAlertRuleRead,

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

func dataSourceSentinelAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, workspaceID.ResourceGroup, "Microsoft.OperationalInsights", workspaceID.WorkspaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Sentinel Alert Rule %q (Workspace: %q) was not found", name, workspaceID)
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule %q (Workspace: %q): %+v", name, workspaceID, err)
	}

	id := alertRuleID(resp.Value)
	if id == nil || *id == "" {
		return fmt.Errorf("nil or empty ID of Sentinel Alert Rule %q (Workspace: %q)", name, workspaceID)
	}

	resId, err := parse.AlertRuleID(*id)
	if err != nil {
		return err
	}
	d.SetId(resId.ID())

	return nil
}
