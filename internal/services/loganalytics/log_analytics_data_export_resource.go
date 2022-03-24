package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsDataExport() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceOperationalinsightsDataExportCreateUpdate,
		Read:   resourceOperationalinsightsDataExportRead,
		Update: resourceOperationalinsightsDataExportCreateUpdate,
		Delete: resourceOperationalinsightsDataExportDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LogAnalyticsDataExportID(id)
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
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.LogAnalyticsDataExportName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"workspace_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceID,
			},

			"destination_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"table_names": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},

			"export_rule_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOperationalinsightsDataExportCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataExportClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := parse.LogAnalyticsWorkspaceID(d.Get("workspace_resource_id").(string))
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	id := parse.NewLogAnalyticsDataExportID(workspace.SubscriptionId, d.Get("resource_group_name").(string), workspace.WorkspaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.DataexportName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_log_analytics_data_export_rule", id.ID())
		}
	}

	parameters := operationalinsights.DataExport{
		DataExportProperties: &operationalinsights.DataExportProperties{
			Destination: &operationalinsights.Destination{
				ResourceID: utils.String(d.Get("destination_resource_id").(string)),
			},
			TableNames: utils.ExpandStringSlice(d.Get("table_names").(*pluginsdk.Set).List()),
			Enable:     utils.Bool(d.Get("enabled").(bool)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.DataexportName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceOperationalinsightsDataExportRead(d, meta)
}

func resourceOperationalinsightsDataExportRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataExportClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsDataExportID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.DataexportName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Log Analytics %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", id.DataexportName, id.ResourceGroup, id.WorkspaceName, err)
	}
	d.Set("name", id.DataexportName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("workspace_resource_id", parse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	if props := resp.DataExportProperties; props != nil {
		d.Set("export_rule_id", props.DataExportID)
		d.Set("destination_resource_id", flattenDataExportDestination(props.Destination))
		d.Set("enabled", props.Enable)
		d.Set("table_names", utils.FlattenStringSlice(props.TableNames))
	}
	return nil
}

func resourceOperationalinsightsDataExportDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataExportClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsDataExportID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.DataexportName); err != nil {
		return fmt.Errorf("deleting Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", id.DataexportName, id.ResourceGroup, id.WorkspaceName, err)
	}
	return nil
}

func flattenDataExportDestination(input *operationalinsights.Destination) string {
	if input == nil {
		return ""
	}

	var resourceID string
	if input.ResourceID != nil {
		resourceID = *input.ResourceID
	}

	return resourceID
}
