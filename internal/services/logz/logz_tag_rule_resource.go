package logz

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var tagRuleName = "default"

func resourceLogzTagRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogzTagRuleCreateUpdate,
		Read:   resourceLogzTagRuleRead,
		Update: resourceLogzTagRuleCreateUpdate,
		Delete: resourceLogzTagRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LogzTagRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"logz_monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogzMonitorName,
			},

			"filtering_tag": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(logz.TagActionInclude),
								string(logz.TagActionExclude),
							}, false),
						},

						"value": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"send_aad_logs": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"send_activity_logs": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"send_subscription_logs": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}
func resourceLogzTagRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewLogzTagRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("logz_monitor_id").(string), tagRuleName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_logz_tag_rule", *existing.ID)
		}
	}

	props := logz.MonitoringTagRules{
		Properties: &logz.MonitoringTagRulesProperties{
			LogRules: expandTagRuleLogRules(d),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName, &props); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogzTagRuleRead(d, meta)
}

func resourceLogzTagRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzTagRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] logz %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("logz_monitor_id", id.MonitorName)
	if props := resp.Properties; props != nil {
		flattenTagRuleLogRules(d, props.LogRules)
	}
	return nil
}

func resourceLogzTagRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzTagRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandTagRuleLogRules(d *pluginsdk.ResourceData) *logz.LogRules {
	return &logz.LogRules{
		SendAadLogs:          utils.Bool(d.Get("send_aad_logs").(bool)),
		SendSubscriptionLogs: utils.Bool(d.Get("send_subscription_logs").(bool)),
		SendActivityLogs:     utils.Bool(d.Get("send_activity_logs").(bool)),
		FilteringTags:        expandTagRuleFilteringTagArray(d.Get("filtering_tag").([]interface{})),
	}
}

func expandTagRuleFilteringTagArray(input []interface{}) *[]logz.FilteringTag {
	results := make([]logz.FilteringTag, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, logz.FilteringTag{
			Name:   utils.String(v["name"].(string)),
			Value:  utils.String(v["value"].(string)),
			Action: logz.TagAction(v["action"].(string)),
		})
	}
	return &results
}

func flattenTagRuleLogRules(d *pluginsdk.ResourceData, input *logz.LogRules) {
	if input == nil {
		return
	}

	d.Set("send_aad_logs", input.SendAadLogs)
	d.Set("send_activity_logs", input.SendActivityLogs)
	d.Set("send_subscription_logs", input.SendSubscriptionLogs)
	d.Set("filtering_tag", flattenTagRuleFilteringTagArray(input.FilteringTags))
}

func flattenTagRuleFilteringTagArray(input *[]logz.FilteringTag) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}
		var action logz.TagAction
		if item.Action != "" {
			action = item.Action
		}
		var value string
		if item.Value != nil {
			value = *item.Value
		}
		results = append(results, map[string]interface{}{
			"name":   name,
			"action": action,
			"value":  value,
		})
	}
	return results
}
