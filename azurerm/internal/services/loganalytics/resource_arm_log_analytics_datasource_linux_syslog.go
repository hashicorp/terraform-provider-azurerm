package loganalytics

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
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

			"severities": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"emerg",
						"alert",
						"crit",
						"err",
						"warning",
						"notice",
						"info",
						"debug",
					}, false),
				},
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
					"local0",
					"local1",
					"local2",
					"local3",
					"local4",
					"local5",
					"local6",
					"local7",
					"lpr",
					"mail",
					"news",
					"syslog",
					"user",
					"uucp",
				}, false),
			},
		},
	}
}

// TODO: We define structure below because of SDK lackes of those definition for now.
//       Once the [issue](https://github.com/Azure/azure-rest-api-specs/issues/9072) addressed,
//       we can switch to using the type directly from SDK.
type dataSourceLinuxSyslog struct {
	Name       string                          `json:"syslogName"`
	Severities []dataSourceLinuxSyslogSeverity `json:"syslogSeverities"`
}

type dataSourceLinuxSyslogSeverity struct {
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
				return fmt.Errorf("checking for existing Log Analytics Data Source Linux Syslog %q (Resource Group %q / Workspace: %q): %+v", name, resourceGroup, workspaceName, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_datasource_linux_syslog", *resp.ID)
		}
	}

	param := operationalinsights.DataSource{
		Kind: operationalinsights.LinuxSyslog,
		Properties: &dataSourceLinuxSyslog{
			Name:       d.Get("syslog_name").(string),
			Severities: expandLogAnalyticsDataSourceLinuxSyslogSeverities(d.Get("severities").(*schema.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, workspaceName, name, param); err != nil {
		return fmt.Errorf("creating Log Analytics DataSource Linux Syslog %q (Resource Group %q / Workspace: %q): %+v", name, resourceGroup, workspaceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, workspaceName, name)
	if err != nil {
		return fmt.Errorf("retrieving Log Analytics Data Source Linux Syslog %q (Resource Group %q / Workspace: %q): %+v", name, resourceGroup, workspaceName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Log Analytics Data Source Linux Syslog %q (Resource Group %q / Workspace: %q) ID", name, resourceGroup, workspaceName)
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
			log.Printf("[DEBUG] Log Analytics Data Source Linux Syslog %q was not found in Resource Group %q in Workspace %q - removing from state!", id.Name, id.ResourceGroup, id.Workspace)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Log Analytics Data Source Linux Syslog %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("workspace_name", id.Workspace)
	if props := resp.Properties; props != nil {
		propStr, err := structure.FlattenJsonToString(props.(map[string]interface{}))
		if err != nil {
			return fmt.Errorf("flattening properties map to json for Log Analytics DataSource Linux Syslog %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
		}

		prop := dataSourceLinuxSyslog{}
		if err := json.Unmarshal([]byte(propStr), &prop); err != nil {
			return fmt.Errorf("decoding properties json for Log Analytics DataSource Linux Syslog %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
		}

		d.Set("syslog_name", prop.Name)
		d.Set("severities", flattenLogAnalyticsDataSourceLinuxSyslogSeverities(prop.Severities))
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
		return fmt.Errorf("deleting Log Analytics Data Source Linux Syslog %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
	}

	return nil
}

func flattenLogAnalyticsDataSourceLinuxSyslogSeverities(severities []dataSourceLinuxSyslogSeverity) []interface{} {
	output := make([]interface{}, 0)
	for _, s := range severities {
		output = append(output, s.Severity)
	}
	return output
}

func expandLogAnalyticsDataSourceLinuxSyslogSeverities(input []interface{}) []dataSourceLinuxSyslogSeverity {
	output := []dataSourceLinuxSyslogSeverity{}
	for _, severity := range input {
		output = append(output, dataSourceLinuxSyslogSeverity{severity.(string)})
	}
	return output
}
