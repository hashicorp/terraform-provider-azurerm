package servicebus

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topicsauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceServiceBusTopicAuthorizationRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusTopicAuthorizationRuleCreateUpdate,
		Read:   resourceServiceBusTopicAuthorizationRuleRead,
		Update: resourceServiceBusTopicAuthorizationRuleCreateUpdate,
		Delete: resourceServiceBusTopicAuthorizationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := topicsauthorizationrule.ParseTopicAuthorizationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: authorizationRuleSchemaFrom(resourceServiceBusTopicAuthorizationRuleSchema()),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(authorizationRuleCustomizeDiff),
	}
}

func resourceServiceBusTopicAuthorizationRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AuthorizationRuleName(),
		},

		//lintignore: S013
		"topic_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: topics.ValidateTopicID,
		},
	}
}

func resourceServiceBusTopicAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.TopicsAuthClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Topic Authorization Rule creation.")

	var id topicsauthorizationrule.TopicAuthorizationRuleId
	if topicIdLit := d.Get("topic_id").(string); topicIdLit != "" {
		topicId, err := topicsauthorizationrule.ParseTopicID(topicIdLit)
		if err != nil {
			return err
		}
		id = topicsauthorizationrule.NewTopicAuthorizationRuleID(topicId.SubscriptionId, topicId.ResourceGroupName, topicId.NamespaceName, topicId.TopicName, d.Get("name").(string))
	}

	if d.IsNewResource() {
		existing, err := client.TopicsGetAuthorizationRule(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_servicebus_topic_authorization_rule", id.ID())
		}
	}

	parameters := topicsauthorizationrule.SBAuthorizationRule{
		Name: utils.String(id.AuthorizationRuleName),
		Properties: &topicsauthorizationrule.SBAuthorizationRuleProperties{
			Rights: *expandTopicAuthorizationRuleRights(d),
		},
	}

	if _, err := client.TopicsCreateOrUpdateAuthorizationRule(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	if err := waitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", id, err)
	}

	return resourceServiceBusTopicAuthorizationRuleRead(d, meta)
}

func resourceServiceBusTopicAuthorizationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.TopicsAuthClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := topicsauthorizationrule.ParseTopicAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.TopicsGetAuthorizationRule(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AuthorizationRuleName)
	d.Set("topic_id", topicsauthorizationrule.NewTopicID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.TopicName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			listen, send, manage := flattenTopicAuthorizationRuleRights(&props.Rights)
			d.Set("listen", listen)
			d.Set("send", send)
			d.Set("manage", manage)
		}
	}

	keysResp, err := client.TopicsListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	if model := keysResp.Model; model != nil {
		d.Set("primary_key", model.PrimaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_key", model.SecondaryKey)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
		d.Set("primary_connection_string_alias", model.AliasPrimaryConnectionString)
		d.Set("secondary_connection_string_alias", model.AliasSecondaryConnectionString)
	}

	return nil
}

func resourceServiceBusTopicAuthorizationRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.TopicsAuthClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := topicsauthorizationrule.ParseTopicAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.TopicsDeleteAuthorizationRule(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)

	if err := waitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for Service Bus Namespace Disaster Recovery Configs (Namespace %q / Resource Group %q): %s", id.NamespaceName, id.ResourceGroupName, err)
	}

	return nil
}

func expandTopicAuthorizationRuleRights(d *pluginsdk.ResourceData) *[]topicsauthorizationrule.AccessRights {
	rights := make([]topicsauthorizationrule.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, topicsauthorizationrule.AccessRightsListen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, topicsauthorizationrule.AccessRightsSend)
	}

	if d.Get("manage").(bool) {
		rights = append(rights, topicsauthorizationrule.AccessRightsManage)
	}

	return &rights
}

func flattenTopicAuthorizationRuleRights(rights *[]topicsauthorizationrule.AccessRights) (listen, send, manage bool) {
	if rights != nil {
		for _, right := range *rights {
			switch right {
			case topicsauthorizationrule.AccessRightsListen:
				listen = true
			case topicsauthorizationrule.AccessRightsSend:
				send = true
			case topicsauthorizationrule.AccessRightsManage:
				manage = true
			default:
				log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
			}
		}
	}

	return listen, send, manage
}
