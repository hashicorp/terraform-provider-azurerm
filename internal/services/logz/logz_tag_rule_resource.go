package logz

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const tagRuleName = "default"

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
			"logz_monitor_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogzMonitorID,
			},

			"tag_filter": {
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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	monitorId, err := parse.LogzMonitorID(d.Get("logz_monitor_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLogzTagRuleID(monitorId.SubscriptionId, monitorId.ResourceGroup, monitorId.MonitorName, tagRuleName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_logz_tag_rule", id.ID())
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

	d.Set("logz_monitor_id", parse.NewLogzMonitorID(id.SubscriptionId, id.ResourceGroup, id.MonitorName).ID())
	if props := resp.Properties; props != nil && props.LogRules != nil {
		d.Set("send_aad_logs", props.LogRules.SendAadLogs)
		d.Set("send_activity_logs", props.LogRules.SendActivityLogs)
		d.Set("send_subscription_logs", props.LogRules.SendSubscriptionLogs)
		if err := d.Set("tag_filter", flattenTagRuleFilteringTagArray(props.LogRules.FilteringTags)); err != nil {
			return fmt.Errorf("setting `tag_filter`: %+v", err)
		}
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
		FilteringTags:        expandTagRuleFilteringTagArray(d.Get("tag_filter").([]interface{})),
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
