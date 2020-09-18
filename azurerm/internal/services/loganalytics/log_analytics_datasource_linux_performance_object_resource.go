package loganalytics

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsDataSourceLinuxPerformanceObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsDataSourceLinuxPerformanceObjectCreateUpdate,
		Read:   resourceArmLogAnalyticsDataSourceLinuxPerformanceObjectRead,
		Update: resourceArmLogAnalyticsDataSourceLinuxPerformanceObjectCreateUpdate,
		Delete: resourceArmLogAnalyticsDataSourceLinuxPerformanceObjectDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.LogAnalyticsDataSourceID(id)
			return err
		}, importLogAnalyticsDataSource(operationalinsights.LinuxPerformanceObject)),

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"workspace_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     ValidateAzureRmLogAnalyticsWorkspaceName,
			},

			"object_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"interval_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(10, 1800),
			},

			"performance_counters": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

type dataSourceLinuxPerformanceObjectProperty struct {
	ObjectName          string                                        `json:"objectName"`
	InstanceName        string                                        `json:"instanceName"`
	IntervalSeconds     int                                           `json:"intervalSeconds"`
	PerformanceCounters []dataSourceLinuxPerformanceObjectCounterName `json:"performanceCounters"`
}

type dataSourceLinuxPerformanceObjectCounterName struct {
	CounterName string `json:"counterName"`
}

func resourceArmLogAnalyticsDataSourceLinuxPerformanceObjectCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	workspaceName := d.Get("workspace_name").(string)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, workspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("failed to check for existing Log Analytics DataSource Linux Performance Object %q (Resource Group %q / Workspace: %q): %+v", name, resourceGroup, workspaceName, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_log_analytics_datasource_linux_performance_object", *resp.ID)
		}
	}

	params := operationalinsights.DataSource{
		Kind: operationalinsights.LinuxPerformanceObject,
		Properties: &dataSourceLinuxPerformanceObjectProperty{
			ObjectName:          d.Get("object_name").(string),
			InstanceName:        d.Get("instance_name").(string),
			IntervalSeconds:     d.Get("interval_seconds").(int),
			PerformanceCounters: expandLogAnalyticsDataSourceLinuxPerformanceObjectCounterName(d.Get("performance_counters").(*schema.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, workspaceName, name, params); err != nil {
		return fmt.Errorf("failed to create Log Analytics DataSource Linux performance object %q (Resource Group %q / Workspace: %q): %+v", name, resourceGroup, workspaceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, workspaceName, name)
	if err != nil {
		return fmt.Errorf("failed to retrieve Log Analytics DataSource Linux performance object %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read ID for Log Analytics DataSource Linux performance object %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmLogAnalyticsDataSourceLinuxPerformanceObjectRead(d, meta)
}

func resourceArmLogAnalyticsDataSourceLinuxPerformanceObjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsDataSourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Log Analytics DataSource Linux performance object %q was not found in Resource Group %q in Workspace %q - removing from state!", id.Name, id.ResourceGroup, id.Workspace)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("failed to retrieve Log Analytics DataSource Linux performance object %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("workspace_name", id.Workspace)
	if props := resp.Properties; props != nil {
		propStr, err := structure.FlattenJsonToString(props.(map[string]interface{}))
		if err != nil {
			return fmt.Errorf("failed to flatten properties map to json for Log Analytics DataSource Linux performance object %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
		}

		prop := &dataSourceLinuxPerformanceObjectProperty{}
		if err := json.Unmarshal([]byte(propStr), &prop); err != nil {
			return fmt.Errorf("failed to decode properties json for Log Analytics DataSource Linux performance object %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
		}

		d.Set("object_name", prop.ObjectName)
		d.Set("instance_name", prop.InstanceName)
		d.Set("interval_seconds", prop.IntervalSeconds)
		d.Set("performance_counters", flattenLogAnalyticsDataSourceLinuxPerformanceObjectCounterName(prop.PerformanceCounters))
	}

	return nil
}

func resourceArmLogAnalyticsDataSourceLinuxPerformanceObjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsDataSourceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
		return fmt.Errorf("failed to delete Log Analytics DataSource Linux Performance Object %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	return nil
}

func expandLogAnalyticsDataSourceLinuxPerformanceObjectCounterName(input []interface{}) []dataSourceLinuxPerformanceObjectCounterName {
	output := []dataSourceLinuxPerformanceObjectCounterName{}
	for _, counterName := range input {
		output = append(output, dataSourceLinuxPerformanceObjectCounterName{counterName.(string)})
	}
	return output
}

func flattenLogAnalyticsDataSourceLinuxPerformanceObjectCounterName(counterNames []dataSourceLinuxPerformanceObjectCounterName) []interface{} {
	output := make([]interface{}, 0)
	for _, e := range counterNames {
		output = append(output, e.CounterName)
	}
	return output
}
