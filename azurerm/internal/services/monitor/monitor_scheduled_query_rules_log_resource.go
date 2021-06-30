package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMonitorScheduledQueryRulesLog() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorScheduledQueryRulesLogCreateUpdate,
		Read:   resourceMonitorScheduledQueryRulesLogRead,
		Update: resourceMonitorScheduledQueryRulesLogCreateUpdate,
		Delete: resourceMonitorScheduledQueryRulesLogDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringDoesNotContainAny("<>*%&:\\?+/"),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"authorized_resource_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"criteria": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dimension": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  "Include",
										ValidateFunc: validation.StringInSlice([]string{
											"Include",
										}, false),
									},
									"values": {
										Type:     pluginsdk.TypeList,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
						"metric_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"data_source_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
			},
			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorScheduledQueryRulesLogCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	action := expandMonitorScheduledQueryRulesLogToMetricAction(d)

	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Monitor Scheduled Query Rules %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_scheduled_query_rules_alert", *existing.ID)
		}
	}

	description := d.Get("description").(string)
	enabledRaw := d.Get("enabled").(bool)

	enabled := insights.True
	if !enabledRaw {
		enabled = insights.False
	}

	location := azure.NormalizeLocation(d.Get("location"))

	source := expandMonitorScheduledQueryRulesCommonSource(d)

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := insights.LogSearchRuleResource{
		Location: utils.String(location),
		LogSearchRule: &insights.LogSearchRule{
			Description: utils.String(description),
			Enabled:     enabled,
			Source:      source,
			Action:      action,
		},
		Tags: expandedTags,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating or updating Scheduled Query Rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Scheduled query rule %q (resource group %q) ID is empty", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceMonitorScheduledQueryRulesLogRead(d, meta)
}

func resourceMonitorScheduledQueryRulesLogRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["scheduledqueryrules"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Scheduled Query Rule %q was not found in Resource Group %q", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Scheduled Query Rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("description", resp.Description)
	if resp.Enabled == insights.True {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}

	action, ok := resp.Action.(insights.LogToMetricAction)
	if !ok {
		return fmt.Errorf("Wrong action type in Scheduled Query Rule %q (resource group %q): %T", name, resourceGroup, resp.Action)
	}
	if err = d.Set("criteria", flattenAzureRmScheduledQueryRulesLogCriteria(action.Criteria)); err != nil {
		return fmt.Errorf("Error setting `criteria`: %+v", err)
	}

	if source := resp.Source; source != nil {
		if source.AuthorizedResources != nil {
			d.Set("authorized_resource_ids", utils.FlattenStringSlice(source.AuthorizedResources))
		}
		if source.DataSourceID != nil {
			d.Set("data_source_id", source.DataSourceID)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMonitorScheduledQueryRulesLogDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["scheduledqueryrules"]

	if resp, err := client.Delete(ctx, resourceGroup, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Scheduled Query Rule %q (resource group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandMonitorScheduledQueryRulesLogCriteria(input []interface{}) *[]insights.Criteria {
	criteria := make([]insights.Criteria, 0)
	if len(input) == 0 {
		return &criteria
	}

	for _, item := range input {
		if item == nil {
			continue
		}
		v, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		dimensions := make([]insights.Dimension, 0)
		for _, dimension := range v["dimension"].(*pluginsdk.Set).List() {
			if dimension == nil {
				continue
			}
			dVal, ok := dimension.(map[string]interface{})
			if !ok {
				continue
			}
			dimensions = append(dimensions, insights.Dimension{
				Name:     utils.String(dVal["name"].(string)),
				Operator: utils.String(dVal["operator"].(string)),
				Values:   utils.ExpandStringSlice(dVal["values"].([]interface{})),
			})
		}

		criteria = append(criteria, insights.Criteria{
			MetricName: utils.String(v["metric_name"].(string)),
			Dimensions: &dimensions,
		})
	}
	return &criteria
}

func expandMonitorScheduledQueryRulesLogToMetricAction(d *pluginsdk.ResourceData) *insights.LogToMetricAction {
	criteriaRaw := d.Get("criteria").([]interface{})
	criteria := expandMonitorScheduledQueryRulesLogCriteria(criteriaRaw)

	action := insights.LogToMetricAction{
		Criteria:  criteria,
		OdataType: insights.OdataTypeMicrosoftWindowsAzureManagementMonitoringAlertsModelsMicrosoftAppInsightsNexusDataContractsResourcesScheduledQueryRulesLogToMetricAction,
	}

	return &action
}
