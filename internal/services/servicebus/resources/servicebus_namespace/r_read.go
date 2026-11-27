// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusNamespaceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ServiceBus.NamespacesClient
	namespaceAuthClient := meta.(*clients.Client).ServiceBus.NamespacesAuthClient
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
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return resourceServiceBusNamespaceFlatten(ctx, d, id, resp.Model, client, namespaceAuthClient)
}

func servicebusTLSVersionDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) (err error) {
	old, new := d.GetChange("minimum_tls_version")
	if old != "" && new == "" {
		err = errors.New("`minimum_tls_version` has been set before, please set a valid value for this property ")
	}
	return
}

func resetNetworkRuleSetForNamespace(ctx context.Context, client *namespaces.NamespacesClient, id namespaces.NamespaceId) error {
	parameters := namespaces.NetworkRuleSet{
		Properties: &namespaces.NetworkRuleSetProperties{
			DefaultAction: pointer.To(namespaces.DefaultActionAllow),
		},
	}

	if _, err := client.CreateOrUpdateNetworkRuleSet(ctx, id, parameters); err != nil {
		return fmt.Errorf("resetting %s: %+v", id, err)
	}

	return nil
}

func networkRuleHash(input any) int {
	v := input.(map[string]any)

	// we are just taking subnet_id into the hash function and ignore the ignore_missing_vnet_service_endpoint to ensure there would be no duplicates of subnet id
	// the service returns this ID with segment resourceGroup and resource group name all in lower cases, to avoid unnecessary diff, we extract this ID and reconstruct this hash code
	return set.HashStringIgnoreCase(v["subnet_id"])
}
