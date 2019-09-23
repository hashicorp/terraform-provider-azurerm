package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMetricAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMetricAlertRuleCreateUpdate,
		Read:   resourceArmMetricAlertRuleRead,
		Update: resourceArmMetricAlertRuleCreateUpdate,
		Delete: resourceArmMetricAlertRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: `The 'azurerm_metric_alertrule' resource is deprecated in favour of the renamed version 'azurerm_monitor_metric_alertrule'.

Information on migrating to the renamed resource can be found here: https://terraform.io/docs/providers/azurerm/guides/migrating-between-renamed-resources.html

As such the existing 'azurerm_metric_alertrule' resource is deprecated and will be removed in the next major version of the AzureRM Provider (2.0).
`,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"operator": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.ConditionOperatorGreaterThan),
					string(insights.ConditionOperatorGreaterThanOrEqual),
					string(insights.ConditionOperatorLessThan),
					string(insights.ConditionOperatorLessThanOrEqual),
				}, true),
			},

			"threshold": {
				Type:     schema.TypeFloat,
				Required: true,
			},

			"period": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"aggregation": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.TimeAggregationOperatorAverage),
					string(insights.TimeAggregationOperatorLast),
					string(insights.TimeAggregationOperatorMaximum),
					string(insights.TimeAggregationOperatorMinimum),
					string(insights.TimeAggregationOperatorTotal),
				}, true),
			},

			"email_action": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_to_service_owners": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},

						"custom_emails": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"webhook_action": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_uri": {
							Type:     schema.TypeString,
							Required: true,
						},

						"properties": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateMetricAlertRuleTags,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmMetricAlertRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.AlertRulesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Alert Rule creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Alert Rule %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_metric_alertrule", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	alertRule, err := expandAzureRmMetricThresholdAlertRule(d)
	if err != nil {
		return err
	}

	alertRuleResource := insights.AlertRuleResource{
		Name:      &name,
		Location:  &location,
		Tags:      tags.Expand(t),
		AlertRule: alertRule,
	}

	if _, err = client.CreateOrUpdate(ctx, resourceGroup, name, alertRuleResource); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM Alert Rule %q (Resource Group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMetricAlertRuleRead(d, meta)
}

func resourceArmMetricAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.AlertRulesClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup, name, err := resourceGroupAndAlertRuleNameFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Metric Alert Rule %q (resource group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Metric Alert Rule %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if alertRule := resp.AlertRule; alertRule != nil {
		d.Set("description", alertRule.Description)
		d.Set("enabled", alertRule.IsEnabled)

		ruleCondition := alertRule.Condition

		if ruleCondition != nil {
			if thresholdRuleCondition, ok := ruleCondition.AsThresholdRuleCondition(); ok && thresholdRuleCondition != nil {
				d.Set("operator", string(thresholdRuleCondition.Operator))
				d.Set("threshold", thresholdRuleCondition.Threshold)
				d.Set("period", thresholdRuleCondition.WindowSize)
				d.Set("aggregation", string(thresholdRuleCondition.TimeAggregation))

				dataSource := thresholdRuleCondition.DataSource

				if dataSource != nil {
					if metricDataSource, ok := dataSource.AsRuleMetricDataSource(); ok && metricDataSource != nil {
						d.Set("resource_id", metricDataSource.ResourceURI)
						d.Set("metric_name", metricDataSource.MetricName)
					}
				}
			}
		}

		email_actions := make([]interface{}, 0)
		webhook_actions := make([]interface{}, 0)

		for _, ruleAction := range *alertRule.Actions {
			if emailAction, ok := ruleAction.AsRuleEmailAction(); ok && emailAction != nil {
				email_action := make(map[string]interface{}, 1)

				if sendToOwners := emailAction.SendToServiceOwners; sendToOwners != nil {
					email_action["send_to_service_owners"] = *sendToOwners
				}

				custom_emails := make([]string, 0)
				if s := emailAction.CustomEmails; s != nil {
					custom_emails = *s
				}
				email_action["custom_emails"] = custom_emails

				email_actions = append(email_actions, email_action)
			} else if webhookAction, ok := ruleAction.AsRuleWebhookAction(); ok && webhookAction != nil {
				webhook_action := make(map[string]interface{}, 1)

				webhook_action["service_uri"] = *webhookAction.ServiceURI

				properties := make(map[string]string)
				for k, v := range webhookAction.Properties {
					if k != "$type" {
						if v != nil {
							properties[k] = *v
						}
					}
				}
				webhook_action["properties"] = properties

				webhook_actions = append(webhook_actions, webhook_action)
			}
		}

		d.Set("email_action", email_actions)
		d.Set("webhook_action", webhook_actions)
	}

	// Return a new tag map filtered by the specified tag names.
	tagMap := tags.Filter(resp.Tags, "$type")

	return tags.FlattenAndSet(d, tagMap)
}

func resourceArmMetricAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.AlertRulesClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup, name, err := resourceGroupAndAlertRuleNameFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error deleting Metric Alert Rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	return err
}

func expandAzureRmMetricThresholdAlertRule(d *schema.ResourceData) (*insights.AlertRule, error) {
	name := d.Get("name").(string)

	resource := d.Get("resource_id").(string)
	metric_name := d.Get("metric_name").(string)

	metricDataSource := insights.RuleMetricDataSource{
		ResourceURI: &resource,
		MetricName:  &metric_name,
	}

	operator := d.Get("operator").(string)
	threshold := d.Get("threshold").(float64)
	period := d.Get("period").(string)
	aggregation := d.Get("aggregation").(string)

	thresholdRuleCondition := insights.ThresholdRuleCondition{
		DataSource:      metricDataSource,
		Operator:        insights.ConditionOperator(operator),
		Threshold:       &threshold,
		TimeAggregation: insights.TimeAggregationOperator(aggregation),
		WindowSize:      &period,
	}

	actions := make([]insights.BasicRuleAction, 0, 2)

	// Email action

	email_actions := d.Get("email_action").([]interface{})

	if len(email_actions) > 0 {
		email_action := email_actions[0].(map[string]interface{})
		emailAction := insights.RuleEmailAction{}

		if v, ok := email_action["custom_emails"]; ok {
			custom_emails := v.([]interface{})

			customEmails := make([]string, 0)
			for _, customEmail := range custom_emails {
				custom_email := customEmail.(string)
				customEmails = append(customEmails, custom_email)
			}

			emailAction.CustomEmails = &customEmails
		}

		if v, ok := email_action["send_to_service_owners"]; ok {
			sendToServiceOwners := v.(bool)
			emailAction.SendToServiceOwners = &sendToServiceOwners
		}

		actions = append(actions, emailAction)
	}

	// Webhook action

	webhook_actions := d.Get("webhook_action").([]interface{})

	if len(webhook_actions) > 0 {
		webhook_action := webhook_actions[0].(map[string]interface{})

		service_uri := webhook_action["service_uri"].(string)

		webhook_properties := make(map[string]*string)

		if v, ok := webhook_action["properties"]; ok {
			properties := v.(map[string]interface{})

			for property_key, property_value := range properties {
				property_string := property_value.(string)
				webhook_properties[property_key] = &property_string
			}
		}

		webhookAction := insights.RuleWebhookAction{
			ServiceURI: &service_uri,
			Properties: webhook_properties,
		}

		actions = append(actions, webhookAction)
	}

	enabled := d.Get("enabled").(bool)

	alertRule := insights.AlertRule{
		Name:      &name,
		Condition: &thresholdRuleCondition,
		Actions:   &actions,
		IsEnabled: &enabled,
	}

	if v, ok := d.GetOk("description"); ok {
		description := v.(string)
		alertRule.Description = &description
	}

	return &alertRule, nil
}

func validateMetricAlertRuleTags(v interface{}, f string) (warnings []string, errors []error) {
	// Normal validation required by any AzureRM resource.
	warnings, errors = tags.Validate(v, f)

	tagsMap := v.(map[string]interface{})

	for k := range tagsMap {
		if strings.EqualFold(k, "$type") {
			errors = append(errors, fmt.Errorf("the %q is not allowed as tag name", k))
		}
	}

	return warnings, errors
}

func resourceGroupAndAlertRuleNameFromId(alertRuleId string) (string, string, error) {
	id, err := azure.ParseAzureResourceID(alertRuleId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["alertrules"]
	resourceGroup := id.ResourceGroup

	return resourceGroup, name, nil
}
