// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logz

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/subaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogzSubAccountTagRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogzSubAccountTagRuleCreate,
		Read:   resourceLogzSubAccountTagRuleRead,
		Update: resourceLogzSubAccountTagRuleUpdate,
		Delete: resourceLogzSubAccountTagRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := tagrules.ParseAccountTagRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"logz_sub_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: subaccount.ValidateAccountID,
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
								string(tagrules.TagActionInclude),
								string(tagrules.TagActionExclude),
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

func resourceLogzSubAccountTagRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subAccountId, err := subaccount.ParseAccountID(d.Get("logz_sub_account_id").(string))
	if err != nil {
		return err
	}

	id := tagrules.NewAccountTagRuleID(subAccountId.SubscriptionId, subAccountId.ResourceGroupName, subAccountId.MonitorName, subAccountId.AccountName, "default")

	existing, err := client.SubAccountTagRulesGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_logz_sub_account_tag_rule", id.ID())
	}

	payload := tagrules.MonitoringTagRules{
		Properties: &tagrules.MonitoringTagRulesProperties{
			LogRules: &tagrules.LogRules{
				FilteringTags:        expandTagRuleFilteringTagArray(d.Get("tag_filter").([]interface{})),
				SendAadLogs:          pointer.To(d.Get("send_aad_logs").(bool)),
				SendSubscriptionLogs: pointer.To(d.Get("send_subscription_logs").(bool)),
				SendActivityLogs:     pointer.To(d.Get("send_activity_logs").(bool)),
			},
		},
	}

	if _, err := client.SubAccountTagRulesCreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogzSubAccountTagRuleRead(d, meta)
}

func resourceLogzSubAccountTagRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseAccountTagRuleID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.SubAccountTagRulesGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := existing.Model

	if d.HasChange("send_aad_logs") {
		payload.Properties.LogRules.SendAadLogs = pointer.To(d.Get("send_aad_logs").(bool))
	}
	if d.HasChange("send_subscription_logs") {
		payload.Properties.LogRules.SendSubscriptionLogs = pointer.To(d.Get("send_subscription_logs").(bool))
	}
	if d.HasChange("send_activity_logs") {
		payload.Properties.LogRules.SendActivityLogs = pointer.To(d.Get("send_activity_logs").(bool))
	}
	if d.HasChange("tag_filter") {
		payload.Properties.LogRules.FilteringTags = expandTagRuleFilteringTagArray(d.Get("tag_filter").([]interface{}))
	}

	if _, err := client.SubAccountTagRulesCreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceLogzSubAccountTagRuleRead(d, meta)
}

func resourceLogzSubAccountTagRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseAccountTagRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SubAccountTagRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("logz_sub_account_id", subaccount.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName, id.AccountName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.LogRules != nil {
			d.Set("send_aad_logs", props.LogRules.SendAadLogs)
			d.Set("send_activity_logs", props.LogRules.SendActivityLogs)
			d.Set("send_subscription_logs", props.LogRules.SendSubscriptionLogs)

			if err := d.Set("tag_filter", flattenTagRuleFilteringTagArray(props.LogRules.FilteringTags)); err != nil {
				return fmt.Errorf("setting `tag_filter`: %+v", err)
			}
		}
	}

	return nil
}

func resourceLogzSubAccountTagRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.TagRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseAccountTagRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.SubAccountTagRulesDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
