package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceLogAnalyticsDataExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmOperationalinsightsDataExportCreateUpdate,
		Read:   resourceArmOperationalinsightsDataExportRead,
		Update: resourceArmOperationalinsightsDataExportCreateUpdate,
		Delete: resourceArmOperationalinsightsDataExportDelete,
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
				ValidateFunc: azure.ValidateResourceID,
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

func resourceArmOperationalinsightsDataExportCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		existing, err := client.Get(ctx, resourceGroup, workspace.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, workspace.Name, err)
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

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, workspace.Name, name, parameters); err != nil {
		return fmt.Errorf("creating/updating Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, workspace.Name, err)
	}

	resp, err := client.Get(ctx, resourceGroup, workspace.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", name, resourceGroup, workspace.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q) ID", name, resourceGroup, workspace.Name)
	}

	d.SetId(*resp.ID)
	return resourceArmOperationalinsightsDataExportRead(d, meta)
}

func resourceArmOperationalinsightsDataExportRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataExportClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsDataExportID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Log Analytics %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", id.Name, id.ResourceGroup, id.WorkspaceName, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("workspace_resource_id", id.WorkspaceID)
	if props := resp.DataExportProperties; props != nil {
		d.Set("export_rule_id", props.DataExportID)
		d.Set("destination_resource_id", flattenArmDataExportDestination(props.Destination))
		d.Set("enabled", props.Enable)
		d.Set("table_names", utils.FlattenStringSlice(props.TableNames))
	}
	return nil
}

func resourceArmOperationalinsightsDataExportDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataExportClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsDataExportID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Log Analytics Data Export Rule %q (Resource Group %q / workspaceName %q): %+v", id.Name, id.ResourceGroup, id.WorkspaceName, err)
	}
	return nil
}

func flattenArmDataExportDestination(input *operationalinsights.Destination) string {
	if input == nil {
		return ""
	}

	var resourceID string
	if input.ResourceID != nil {
		resourceID = *input.ResourceID
	}

	return resourceID
}
