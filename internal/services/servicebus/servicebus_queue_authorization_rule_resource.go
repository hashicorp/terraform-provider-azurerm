package servicebus

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queuesauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceServiceBusQueueAuthorizationRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusQueueAuthorizationRuleCreateUpdate,
		Read:   resourceServiceBusQueueAuthorizationRuleRead,
		Update: resourceServiceBusQueueAuthorizationRuleCreateUpdate,
		Delete: resourceServiceBusQueueAuthorizationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := queuesauthorizationrule.ParseQueueAuthorizationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceServiceBusqueueAuthorizationRuleSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(authorizationRuleCustomizeDiff),
	}
}

func resourceServiceBusqueueAuthorizationRuleSchema() map[string]*pluginsdk.Schema {
	return authorizationRuleSchemaFrom(map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AuthorizationRuleName(),
		},

		//lintignore: S013
		"queue_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: queues.ValidateQueueID,
		},
	})
}

func resourceServiceBusQueueAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesAuthClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	queueId, err := queuesauthorizationrule.ParseQueueID(d.Get("queue_id").(string))
	if err != nil {
		return err
	}

	id := queuesauthorizationrule.NewQueueAuthorizationRuleID(queueId.SubscriptionId, queueId.ResourceGroupName, queueId.NamespaceName, queueId.QueueName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.QueuesGetAuthorizationRule(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_servicebus_queue_authorization_rule", id.ID())
		}
	}

	parameters := queuesauthorizationrule.SBAuthorizationRule{
		Name: utils.String(id.AuthorizationRuleName),
		Properties: &queuesauthorizationrule.SBAuthorizationRuleProperties{
			Rights: *expandQueueAuthorizationRuleRights(d),
		},
	}

	if _, err := client.QueuesCreateOrUpdateAuthorizationRule(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	if err := waitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", id, err)
	}

	return resourceServiceBusQueueAuthorizationRuleRead(d, meta)
}

func resourceServiceBusQueueAuthorizationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesAuthClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := queuesauthorizationrule.ParseQueueAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.QueuesGetAuthorizationRule(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AuthorizationRuleName)
	d.Set("queue_id", queuesauthorizationrule.NewQueueID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.QueueName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			listen, send, manage := flattenQueueAuthorizationRuleRights(&props.Rights)
			d.Set("manage", manage)
			d.Set("listen", listen)
			d.Set("send", send)
		}
	}

	keysResp, err := client.QueuesListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	if keysModel := keysResp.Model; keysModel != nil {
		d.Set("primary_key", keysModel.PrimaryKey)
		d.Set("primary_connection_string", keysModel.PrimaryConnectionString)
		d.Set("secondary_key", keysModel.SecondaryKey)
		d.Set("secondary_connection_string", keysModel.SecondaryConnectionString)
		d.Set("primary_connection_string_alias", keysModel.AliasPrimaryConnectionString)
		d.Set("secondary_connection_string_alias", keysModel.AliasSecondaryConnectionString)
	}

	return nil
}

func resourceServiceBusQueueAuthorizationRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.QueuesAuthClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := queuesauthorizationrule.ParseQueueAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.QueuesDeleteAuthorizationRule(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	if err := waitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", *id, err)
	}

	return nil
}

func expandQueueAuthorizationRuleRights(d *pluginsdk.ResourceData) *[]queuesauthorizationrule.AccessRights {
	rights := make([]queuesauthorizationrule.AccessRights, 0)

	if d.Get("listen").(bool) {
		rights = append(rights, queuesauthorizationrule.AccessRightsListen)
	}

	if d.Get("send").(bool) {
		rights = append(rights, queuesauthorizationrule.AccessRightsSend)
	}

	if d.Get("manage").(bool) {
		rights = append(rights, queuesauthorizationrule.AccessRightsManage)
	}

	return &rights
}

func flattenQueueAuthorizationRuleRights(rights *[]queuesauthorizationrule.AccessRights) (listen, send, manage bool) {
	// zero (initial) value for a bool in go is false

	if rights != nil {
		for _, right := range *rights {
			switch right {
			case queuesauthorizationrule.AccessRightsListen:
				listen = true
			case queuesauthorizationrule.AccessRightsSend:
				send = true
			case queuesauthorizationrule.AccessRightsManage:
				manage = true
			default:
				log.Printf("[DEBUG] Unknown Authorization Rule Right '%s'", right)
			}
		}
	}

	return listen, send, manage
}
