package loganalytics

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/state"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"time"
)

func resourceArmLogAnalyticsDataSourceLinuxSyslog() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsDataSourceLinuxSyslogCreateUpdate,
		Read:   resourceArmLogAnalyticsDataSourceLinuxSyslogRead,
		Update: resourceArmLogAnalyticsDataSourceLinuxSyslogCreateUpdate,
		Delete: resourceArmLogAnalyticsDataSourceLinuxSyslogDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.LogAnalyticsDataSourceID(id)
			return err
		}, importLogAnalyticsDataSource(operationalinsights.LinuxSyslog)),

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

			"syslog_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"auth",
					"authpriv",
					"cron",
					"daemon",
					"ftp",
					"kern",
					"mail",
					"syslog",
					"user",
					"uucp",
				}, false),
			},

			"syslog_severities": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Set:      set.HashStringIgnoreCase,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					// API backend accepts event_types case-insensitively
					ValidateFunc:     validation.StringInSlice([]string{"emerg", "alert", "crit", "err", "warning", "notice", "info", "debug"}, true),
					StateFunc:        state.IgnoreCase,
					DiffSuppressFunc: suppress.CaseDifference,
				},
			},
		},
	}
}

type dataSourceLinuxSysLogProperty struct {
	SysLogName       string                          `json:"syslogName"`
	SysLogSeverities []dataSourceLinuxSysLogSeverity `json:"syslogSeverities"`
}

type dataSourceLinuxSysLogSeverity struct {
	Severity string `json:"severity"`
}

func resourceArmLogAnalyticsDataSourceLinuxSyslogCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("failed to check for existing Log Analytics DataSource Linux syslog %q (Resource Group %q / Workspace: %q): %+v", name, resourceGroup, workspaceName, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_log_analytics_datasource_linux_syslog", *resp.ID)
		}
	}

	params := operationalinsights.DataSource{
		Kind: operationalinsights.LinuxSyslog,
		Properties: &dataSourceLinuxSysLogProperty{
			SysLogName:       d.Get("syslog_name").(string),
			SysLogSeverities: expandLogAnalyticsDataSourceLinuxSyslogSeverity(d.Get("syslog_severities").(*schema.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, workspaceName, name, params); err != nil {
		return fmt.Errorf("failed to create Log Analytics DataSource Linux syslog %q (Resource Group %q / Workspace: %q): %+v", name, resourceGroup, workspaceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, workspaceName, name)
	if err != nil {
		return fmt.Errorf("failed to retrieve Log Analytics DataSource Linux syslog %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read ID for Log Analytics DataSource Linux Syslog %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmLogAnalyticsDataSourceLinuxSyslogRead(d, meta)
}

func resourceArmLogAnalyticsDataSourceLinuxSyslogRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Log Analytics DataSource Linux syslog %q was not found in Resource Group %q in Workspace %q - removing from state!", id.Name, id.ResourceGroup, id.Workspace)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("failed to retrieve Log Analytics DataSource Linux syslog %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("workspace_name", id.Workspace)
	if props := resp.Properties; props != nil {
		propStr, err := structure.FlattenJsonToString(props.(map[string]interface{}))
		if err != nil {
			return fmt.Errorf("failed to flatten properties map to json for Log Analytics DataSource Linux syslog %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
		}

		prop := &dataSourceLinuxSysLogProperty{}
		if err := json.Unmarshal([]byte(propStr), &prop); err != nil {
			return fmt.Errorf("failed to decode properties json for Log Analytics DataSource Linux syslog %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
		}

		d.Set("syslog_name", prop.SysLogName)
		d.Set("syslog_severities", flattenLogAnalyticsDataSourceLinuxSyslogSeverity(prop.SysLogSeverities))
	}

	return nil
}

func resourceArmLogAnalyticsDataSourceLinuxSyslogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogAnalyticsDataSourceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Workspace, id.Name); err != nil {
		return fmt.Errorf("failed to delete Log Analytics DataSource Linux syslog %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	return nil
}

func expandLogAnalyticsDataSourceLinuxSyslogSeverity(input []interface{}) []dataSourceLinuxSysLogSeverity {
	output := []dataSourceLinuxSysLogSeverity{}
	for _, severity := range input {
		output = append(output, dataSourceLinuxSysLogSeverity{severity.(string)})
	}
	return output
}

func flattenLogAnalyticsDataSourceLinuxSyslogSeverity(severities []dataSourceLinuxSysLogSeverity) []interface{} {
	output := make([]interface{}, 0)
	for _, e := range severities {
		output = append(output, e.Severity)
	}
	return output
}
