// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceServiceBusNamespaceAuthorizationRule() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
		Read: dataSourceServiceBusNamespaceAuthorizationRuleRead,

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

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
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
		},
	}

	if !features.FivePointOh() {
		r.Schema["namespace_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.NamespaceName,
			AtLeastOneOf: []string{"namespace_id", "resource_group_name", "namespace_name"},
			Deprecated:   "`namespace_name` will be removed in favour of the property `namespace_id` in v5.0 of the AzureRM Provider.",
		}

		r.Schema["resource_group_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: resourcegroups.ValidateName,
			AtLeastOneOf: []string{"namespace_id", "resource_group_name", "namespace_name"},
			Deprecated:   "`resource_group_name` will be removed in favour of the property `namespace_id` in v5.0 of the AzureRM Provider.",
		}

		r.Schema["namespace_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: namespaces.ValidateNamespaceID,
			AtLeastOneOf: []string{"namespace_id", "resource_group_name", "namespace_name"},
		}
	}

	return r
}

func dataSourceServiceBusNamespaceAuthorizationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesAuthClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var id namespacesauthorizationrule.AuthorizationRuleId // remove during `features.FivePointOh()` clean up
	if features.FivePointOh() {
		namespaceId, err := namespacesauthorizationrule.ParseNamespaceID(d.Get("namespace_id").(string))
		if err != nil {
			return err
		}

		id = namespacesauthorizationrule.NewAuthorizationRuleID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))
	} else {
		var resourceGroup string
		var namespaceName string
		if v, ok := d.Get("namespace_id").(string); ok && v != "" {
			namespaceId, err := namespacesauthorizationrule.ParseNamespaceID(v)
			if err != nil {
				return err
			}
			resourceGroup = namespaceId.ResourceGroupName
			namespaceName = namespaceId.NamespaceName
		} else {
			resourceGroup = d.Get("resource_group_name").(string)
			namespaceName = d.Get("namespace_name").(string)
		}

		id = namespacesauthorizationrule.NewAuthorizationRuleID(subscriptionId, resourceGroup, namespaceName, d.Get("name").(string))
	}

	resp, err := client.NamespacesGetAuthorizationRule(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	keysResp, err := client.NamespacesListKeys(ctx, id)
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

	d.SetId(id.ID())

	return nil
}
