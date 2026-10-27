// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRelayNamespaceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Relay.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	authRuleId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, "RootManageSharedAccessKey")
	keysResp, err := client.ListKeys(ctx, authRuleId)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	d.Set("name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
		}

		if props := model.Properties; props != nil {
			d.Set("metric_id", props.MetricId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	if model := keysResp.Model; model != nil {
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("primary_key", model.PrimaryKey)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
		d.Set("secondary_key", model.SecondaryKey)
	}

	return nil
}
