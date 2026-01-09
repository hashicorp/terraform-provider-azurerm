// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/disasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceServiceBusNamespaceDisasterRecoveryConfig() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceServiceBusNamespaceDisasterRecoveryConfigRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"namespace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: namespaces.ValidateNamespaceID,
			},

			"alias_authorization_rule_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"partner_namespace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}

	if !features.FivePointOh() {
		resource.Schema["namespace_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
			AtLeastOneOf: []string{"resource_group_name", "namespace_name", "namespace_id"},
		}
		resource.Schema["namespace_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.NamespaceName,
			AtLeastOneOf: []string{"resource_group_name", "namespace_name", "namespace_id"},
			Deprecated:   "`namespace_name` will be removed in favour of the property `namespace_id` in version 5.0 of the AzureRM Provider.",
		}
		resource.Schema["resource_group_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: resourcegroups.ValidateName,
			AtLeastOneOf: []string{"resource_group_name", "namespace_name", "namespace_id"},
			Deprecated:   "`resource_group_name` will be removed in favour of the property `namespace_id` in version 5.0 of the AzureRM Provider.",
		}
	}

	return resource
}

func dataSourceServiceBusNamespaceDisasterRecoveryConfigRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var resourceGroup string
	var namespaceName string
	if v, ok := d.Get("namespace_id").(string); ok && v != "" {
		namespaceId, err := disasterrecoveryconfigs.ParseNamespaceID(v)
		if err != nil {
			return fmt.Errorf("parsing namespace ID %q: %+v", v, err)
		}
		resourceGroup = namespaceId.ResourceGroupName
		namespaceName = namespaceId.NamespaceName

		if !features.FivePointOh() && namespaceId.NamespaceName == "" {
			resourceGroup = d.Get("resource_group_name").(string)
			namespaceName = d.Get("namespace_name").(string)
		}
	}

	id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID(subscriptionId, resourceGroup, namespaceName, d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.DisasterRecoveryConfigName)

	if !features.FivePointOh() {
		d.Set("resource_group_name", id.ResourceGroupName)
		d.Set("namespace_name", id.NamespaceName)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("partner_namespace_id", props.PartnerNamespace)
		}
	}

	d.SetId(id.ID())

	// the auth rule cannot be retrieved by dr config name, the shared access policy should either be specified by user or using the default one which is `RootManageSharedAccessKey`
	authRuleId := disasterrecoveryconfigs.NewDisasterRecoveryConfigAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.DisasterRecoveryConfigName, serviceBusNamespaceDefaultAuthorizationRule)
	if input := d.Get("alias_authorization_rule_id").(string); input != "" {
		ruleId, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigAuthorizationRuleID(input)
		if err != nil {
			return fmt.Errorf("parsing primary namespace auth rule id error: %+v", err)
		}
		authRuleId = *ruleId
	}

	keys, err := client.ListKeys(ctx, authRuleId)
	if err != nil {
		log.Printf("[WARN] listing default keys for %s: %+v", id, err)
	} else {
		if keysModel := keys.Model; keysModel != nil {
			d.Set("primary_connection_string_alias", keysModel.AliasPrimaryConnectionString)
			d.Set("secondary_connection_string_alias", keysModel.AliasSecondaryConnectionString)
			d.Set("default_primary_key", keysModel.PrimaryKey)
			d.Set("default_secondary_key", keysModel.SecondaryKey)
		}
	}
	return nil
}
