package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/arm/monitor"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
)

func resourceArmMetricAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMetricAlertRuleCreateOrUpdate,
		Read:   resourceArmMetricAlertRuleRead,
		Update: resourceArmMetricAlertRuleCreateOrUpdate,
		Delete: resourceArmMetricAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

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

			"resource": {
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
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(monitor.ConditionOperatorGreaterThan),
					string(monitor.ConditionOperatorGreaterThanOrEqual),
					string(monitor.ConditionOperatorLessThan),
					string(monitor.ConditionOperatorLessThanOrEqual),
				}, true),
			},

			"threshold": {
				Type:     schema.TypeFloat,
				Required: true,
			},

			"period": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIso8601Duration(),
			},

			"aggregation": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(monitor.TimeAggregationOperatorAverage),
					string(monitor.TimeAggregationOperatorLast),
					string(monitor.TimeAggregationOperatorMaximum),
					string(monitor.TimeAggregationOperatorMinimum),
					string(monitor.TimeAggregationOperatorTotal),
				}, true),
			},

			"email_action": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_owners": {
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmMetricAlertRuleCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorAlertRulesClient

	log.Printf("[INFO] preparing arguments for AzureRM Alert Rule creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	alertRule, _ := expandAzureRmMetricThresholdAlertRule(d)

	alertRuleResource := monitor.AlertRuleResource{
		Name:      &name,
		Location:  &location,
		Tags:      expandTags(tags),
		AlertRule: alertRule,
	}

	_, err := client.CreateOrUpdate(resGroup, name, alertRuleResource)
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM Alert Rule '%s' (Resource Group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMetricAlertRuleRead(d, meta)
}

func resourceArmMetricAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorAlertRulesClient

	resGroup, name, err := resourceGroupAndAlertRuleNameFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Metric Alert Rule %q (resource group %q) was not found - removing from state", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Metric Alert Rule '%s': %s", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	alertRule := resp.AlertRule

	if alertRule != nil {
		d.Set("description", *alertRule.Description)
		d.Set("enabled", *alertRule.IsEnabled)

		ruleCondition := alertRule.Condition

		if ruleCondition != nil {
			thresholdRuleCondition, _ := ruleCondition.AsThresholdRuleCondition()

			d.Set("operator", string(thresholdRuleCondition.Operator))
			d.Set("threshold", float64(*thresholdRuleCondition.Threshold))
			d.Set("period", *thresholdRuleCondition.WindowSize)
			d.Set("aggregation", string(thresholdRuleCondition.TimeAggregation))

			dataSource := thresholdRuleCondition.DataSource

			if dataSource != nil {
				metricDataSource, _ := dataSource.AsRuleMetricDataSource()

				d.Set("resource", *metricDataSource.ResourceURI)
				d.Set("metric_name", *metricDataSource.MetricName)
			}
		}

		email_actions := make([]interface{}, 0)
		webhook_actions := make([]interface{}, 0)

		for _, ruleAction := range *alertRule.Actions {
			if emailAction, ok := ruleAction.AsRuleEmailAction(); ok {
				email_action := make(map[string]interface{}, 1)

				email_action["service_owners"] = *emailAction.SendToServiceOwners

				custom_emails := []string{}
				for _, custom_email := range *emailAction.CustomEmails {
					custom_emails = append(custom_emails, custom_email)
				}

				email_action["custom_emails"] = custom_emails

				email_actions = append(email_actions, email_action)
			} else if webhookAction, ok := ruleAction.AsRuleWebhookAction(); ok {
				webhook_action := make(map[string]interface{}, 1)

				webhook_action["service_uri"] = *webhookAction.ServiceURI

				properties := make(map[string]string, len(*webhookAction.Properties))
				for k, v := range *webhookAction.Properties {
					if k != "$type" {
						properties[k] = *v
					}
				}
				webhook_action["properties"] = properties

				webhook_actions = append(webhook_actions, webhook_action)
			}
		}

		d.Set("email_action", email_actions)
		d.Set("webhook_action", webhook_actions)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmMetricAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorAlertRulesClient

	resGroup, name, err := resourceGroupAndAlertRuleNameFromId(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(resGroup, name)

	return err
}

func expandAzureRmMetricThresholdAlertRule(d *schema.ResourceData) (*monitor.AlertRule, error) {
	name := d.Get("name").(string)

	resource := d.Get("resource").(string)
	metric_name := d.Get("metric_name").(string)

	metricDataSource := monitor.RuleMetricDataSource{
		ResourceURI: &resource,
		MetricName:  &metric_name,
	}

	operator := d.Get("operator").(string)
	threshold := d.Get("threshold").(float64)
	period := d.Get("period").(string)
	aggregation := d.Get("aggregation").(string)

	thresholdRuleCondition := monitor.ThresholdRuleCondition{
		DataSource:      metricDataSource,
		Operator:        monitor.ConditionOperator(operator),
		Threshold:       &threshold,
		TimeAggregation: monitor.TimeAggregationOperator(aggregation),
		WindowSize:      &period,
	}

	actions := make([]monitor.RuleAction, 0, 2)

	// Email action

	email_actions := d.Get("email_action").([]interface{})

	if len(email_actions) > 0 {
		email_action := email_actions[0].(map[string]interface{})
		emailAction := monitor.RuleEmailAction{}

		if v, ok := email_action["custom_emails"]; ok {
			custom_emails := v.([]interface{})

			customEmails := make([]string, 0)
			for _, customEmail := range custom_emails {
				custom_email := customEmail.(string)
				customEmails = append(customEmails, custom_email)
			}

			emailAction.CustomEmails = &customEmails
		}

		if v, ok := email_action["service_owners"]; ok {
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

		webhookAction := monitor.RuleWebhookAction{
			ServiceURI: &service_uri,
			Properties: &webhook_properties,
		}

		actions = append(actions, webhookAction)
	}

	enabled := d.Get("enabled").(bool)

	alertRule := monitor.AlertRule{
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

func resourceGroupAndAlertRuleNameFromId(alertRuleId string) (string, string, error) {
	id, err := parseAzureResourceID(alertRuleId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["alertrules"]
	resGroup := id.ResourceGroup

	return resGroup, name, nil
}
