package sentinel

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceSentinelAlertRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSentinelAlertRuleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},
		},
	}
}

func dataSourceSentinelAlertRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := workspaces.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	resp, err := client.Get(ctx, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Sentinel Alert Rule %q was not found", id)
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return nil
}
