package loganalytics

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsWorkspaceTable() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsWorkspaceTableUpdate,
		Read:   resourceLogAnalyticsWorkspaceTableRead,
		Update: resourceLogAnalyticsWorkspaceTableUpdate,
		Delete: resourceLogAnalyticsWorkspaceTableSetDefaults,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LogAnalyticsWorkspaceTableID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WorkspaceV0ToV1{},
			1: migration.WorkspaceV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"retention_in_days": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.Any(validation.IntBetween(30, 730), validation.IntInSlice([]int{7})),
			},
		},
	}
}

func resourceLogAnalyticsWorkspaceTableUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	tablesClient := meta.(*clients.Client).LogAnalytics.TablesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tableName := d.Get("name").(string)
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace Table %s update.", tableName)

	workspaceId, err := parse.LogAnalyticsWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("invalid workspace object ID for table %s: %s", tableName, err)
	}

	retentionInDays := int32(d.Get("retention_in_days").(int))
	result, err := tablesClient.Update(ctx, workspaceId.ResourceGroup, workspaceId.WorkspaceName, tableName, operationalinsights.Table{
		TableProperties: &operationalinsights.TableProperties{
			RetentionInDays: &retentionInDays,
		},
	})
	if result.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update table %s in workspace %s in resource group %s: %s", tableName, workspaceId.WorkspaceName, workspaceId.ResourceGroup, result.Status)
	}
	if err != nil {
		return fmt.Errorf("failed to update table %s in workspace %s in resource group %s: %s", tableName, workspaceId.WorkspaceName, workspaceId.ResourceGroup, err)
	}

	d.SetId(*result.ID)

	return nil
}

func resourceLogAnalyticsWorkspaceTableRead(d *pluginsdk.ResourceData, meta interface{}) error {
	tablesClient := meta.(*clients.Client).LogAnalytics.TablesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsWorkspaceTableID(d.Id())
	if err != nil {
		return err
	}

	result, err := tablesClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.TableName)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Log Analytics workspace '%s' table '%s': %+v", id.WorkspaceName, id.TableName, err)
	}

	d.Set("name", result.Name)
	d.Set("workspace_id", parse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	d.Set("retention_in_days", result.TableProperties.RetentionInDays)

	d.SetId(*result.ID)

	return nil
}

func resourceLogAnalyticsWorkspaceTableSetDefaults(d *pluginsdk.ResourceData, meta interface{}) error {
	tablesClient := meta.(*clients.Client).LogAnalytics.TablesClient
	workspacesClient := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tableName := d.Get("name").(string)
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace Table %s restore defaults.", tableName)

	workspaceId, err := parse.LogAnalyticsWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("invalid workspace object ID for table %s: %s", tableName, err)
	}

	// setting the table's RetentionInDays to nil in the update call does not seem to set the table's RetentionInDays
	// to the workspace's default so read RetentionInDays from the workspace and set the table's RetentionInDays
	// to that.
	workspaceResult, err := workspacesClient.Get(ctx, workspaceId.ResourceGroup, workspaceId.WorkspaceName)
	if workspaceResult.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to read workspace while restoring table %s in workspace %s in resource group %s to the default workspace retention: %s",
			tableName, workspaceId.WorkspaceName, workspaceId.ResourceGroup, workspaceResult.Status)
	}
	if err != nil {
		return fmt.Errorf(
			"failed to read workspace while restoring table %s in workspace %s in resource group %s to the default workspace retention: %s",
			tableName, workspaceId.WorkspaceName, workspaceId.ResourceGroup, err)
	}

	log.Printf("[INFO] setting AzureRM Log Analytics Workspace Table %s to default retention %d.", tableName, *workspaceResult.RetentionInDays)
	tableResult, err := tablesClient.Update(ctx, workspaceId.ResourceGroup, workspaceId.WorkspaceName, tableName, operationalinsights.Table{
		TableProperties: &operationalinsights.TableProperties{
			RetentionInDays: workspaceResult.RetentionInDays,
		},
	})
	if tableResult.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to restore table %s in workspace %s in resource group %s to the default workspace retention: %s",
			tableName, workspaceId.WorkspaceName, workspaceId.ResourceGroup, tableResult.Status)
	}
	if err != nil {
		return fmt.Errorf(
			"failed to restore table %s in workspace %s in resource group %s to the default workspace retention: %s",
			tableName, workspaceId.WorkspaceName, workspaceId.ResourceGroup, err)
	}

	d.Set("retention_in_days", workspaceResult.RetentionInDays)
	d.SetId(*tableResult.ID)

	return nil
}
