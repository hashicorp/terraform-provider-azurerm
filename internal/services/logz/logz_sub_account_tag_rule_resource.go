package logz

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logz/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogzSubAccountTagRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogzSubAccountTagRuleCreateUpdate,
		Read:   resourceLogzSubAccountTagRuleRead,
		Update: resourceLogzSubAccountTagRuleCreateUpdate,
		Delete: resourceLogzSubAccountTagRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LogzSubAccountTagRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"logz_sub_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogzSubAccountID,
			},

			"tag_filter": schemaTagFilter(),

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
func resourceLogzSubAccountTagRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountTagRuleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subAccountId, err := parse.LogzSubAccountID(d.Get("logz_sub_account_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLogzSubAccountTagRuleID(subAccountId.SubscriptionId, subAccountId.ResourceGroup, subAccountId.MonitorName, subAccountId.AccountName, TagRuleName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.AccountName, id.TagRuleName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_logz_sub_account_tag_rule", *existing.ID)
		}
	}

	props := logz.MonitoringTagRules{
		Properties: &logz.MonitoringTagRulesProperties{
			LogRules: expandTagRuleLogRules(d),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.MonitorName, id.AccountName, id.TagRuleName, &props); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogzSubAccountTagRuleRead(d, meta)
}

func resourceLogzSubAccountTagRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountTagRuleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzSubAccountTagRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MonitorName, id.AccountName, id.TagRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] logz %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("logz_sub_account_id", parse.NewLogzSubAccountID(id.SubscriptionId, id.ResourceGroup, id.MonitorName, id.AccountName).ID())
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

func resourceLogzSubAccountTagRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logz.SubAccountTagRuleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogzSubAccountTagRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.MonitorName, id.AccountName, id.TagRuleName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
