package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceLogAnalyticsDataExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceOperationalinsightsDataExportCreateUpdate,
		Read:   resourceOperationalinsightsDataExportRead,
		Update: resourceOperationalinsightsDataExportCreateUpdate,
		Delete: resourceOperationalinsightsDataExportDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.LogAnalyticsDataExportName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"workspace_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceID,
			},

			"destination_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"table_names": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},

			"export_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOperationalinsightsDataExportCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataExportClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	workspace, err := parse.LogAnalyticsWorkspaceID(d.Get("workspace_resource_id").(string))
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, workspace.WorkspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, workspace.WorkspaceName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_data_export_rule", *existing.ID)
		}
	}

	parameters := operationalinsights.DataExport{
		DataExportProperties: &operationalinsights.DataExportProperties{
			Destination: &operationalinsights.Destination{
				ResourceID: utils.String(d.Get("destination_resource_id").(string)),
			},
			TableNames: utils.ExpandStringSlice(d.Get("table_names").(*schema.Set).List()),
			Enable:     utils.Bool(d.Get("enabled").(bool)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, workspace.WorkspaceName, name, parameters); err != nil {
		return fmt.Errorf("creating/updating Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, workspace.WorkspaceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, workspace.WorkspaceName, name)
	if err != nil {
		return fmt.Errorf("retrieving Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, workspace.WorkspaceName, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q) ID", name, resourceGroup, workspace.WorkspaceName)
	}

	d.SetId(*resp.ID)
	return resourceOperationalinsightsDataExportRead(d, meta)
}

func resourceOperationalinsightsDataExportRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceOperationalinsightsDataExportDelete(d *schema.ResourceData, meta interface{}) error {
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
