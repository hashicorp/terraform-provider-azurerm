package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceServiceBusNamespaceAuthorizationRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusNamespaceAuthorizationRuleCreateUpdate,
		Read:   resourceServiceBusNamespaceAuthorizationRuleRead,
		Update: resourceServiceBusNamespaceAuthorizationRuleCreateUpdate,
		Delete: resourceServiceBusNamespaceAuthorizationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := namespacesauthorizationrule.ParseAuthorizationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		// function takes a schema map and adds the authorization rule properties to it
		Schema: resourceServiceBusNamespaceAuthorizationRuleSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(authorizationRuleCustomizeDiff),
	}
}

func resourceServiceBusNamespaceAuthorizationRuleSchema() map[string]*pluginsdk.Schema {
	return authorizationRuleSchemaFrom(map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AuthorizationRuleName(),
		},
		//lintignore: S013
		"namespace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
		},
	})
}

func resourceServiceBusNamespaceAuthorizationRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesAuthClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	namespaceAuthId, err := namespacesauthorizationrule.ParseNamespaceID(d.Get("namespace_id").(string))
	if err != nil {
		return err
	}

	id := namespacesauthorizationrule.NewAuthorizationRuleID(namespaceAuthId.SubscriptionId, namespaceAuthId.ResourceGroupName, namespaceAuthId.NamespaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.NamespacesGetAuthorizationRule(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace_authorization_rule", id.ID())
		}
	}

	parameters := namespacesauthorizationrule.SBAuthorizationRule{
		Name: utils.String(id.AuthorizationRuleName),
		Properties: &namespacesauthorizationrule.SBAuthorizationRuleProperties{
			Rights: *expandAuthorizationRuleRights(d),
		},
	}

	if _, err := client.NamespacesCreateOrUpdateAuthorizationRule(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	if err := waitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", id, err)
	}

	return resourceServiceBusNamespaceAuthorizationRuleRead(d, meta)
}

func resourceServiceBusNamespaceAuthorizationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesAuthClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespacesauthorizationrule.ParseAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.NamespacesGetAuthorizationRule(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AuthorizationRuleName)
	d.Set("namespace_id", namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			listen, send, manage := flattenAuthorizationRuleRights(&props.Rights)
			d.Set("manage", manage)
			d.Set("listen", listen)
			d.Set("send", send)
		}
	}

	keysResp, err := client.NamespacesListKeys(ctx, *id)
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

func resourceServiceBusNamespaceAuthorizationRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesAuthClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespacesauthorizationrule.ParseAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.NamespacesDeleteAuthorizationRule(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	if err := waitForPairedNamespaceReplication(ctx, meta, namespaceId, d.Timeout(pluginsdk.TimeoutUpdate)); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", *id, err)
	}

	return nil
}
