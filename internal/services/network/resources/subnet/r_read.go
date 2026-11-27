// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subnet

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetRead(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).Subnets

	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if err := resourceSubnetFlatten(d, *id, resp.Model); err != nil {
		return fmt.Errorf("encoding %s: %+v", id, err)
	}

	return nil
}

// subnetServiceEndpointDiffSuppress is a TEMPORARY solution to suppress diffs for a service endpoint list where the
// individual entries have no change, and only the order of endpoints is different.
// The Azure API will not persist an order change, which results in an infinite diff loop if order was the only change.
// This function should be removed when the SDK v2 supports this feature for TypeList natively.
func subnetServiceEndpointDiffSuppress(_, _, _ string, d *schema.ResourceData) bool {
	if d == nil {
		return false
	}

	listKey := "service_endpoint"

	oldRaw, newRaw := d.GetChange(listKey)
	if oldRaw == nil || newRaw == nil {
		return false
	}

	oldList, okOld := oldRaw.([]any)
	newList, okNew := newRaw.([]any)
	if !okOld || !okNew {
		return false
	}

	if len(oldList) != len(newList) {
		return false
	}

	type endpointItem struct {
		service           string
		networkIdentifier string
	}

	canonicalize := func(rawList []any) []endpointItem {
		items := make([]endpointItem, 0, len(rawList))
		for _, raw := range rawList {
			m, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			item := endpointItem{}
			if s, ok := m["service"].(string); ok {
				item.service = s
			}
			if ni, ok := m["network_identifier"].(string); ok {
				item.networkIdentifier = ni
			}
			items = append(items, item)
		}

		sort.Slice(items, func(i, j int) bool {
			if items[i].service != items[j].service {
				return items[i].service < items[j].service
			}
			return items[i].networkIdentifier < items[j].networkIdentifier
		})

		return items
	}

	canonicalOld := canonicalize(oldList)
	canonicalNew := canonicalize(newList)

	return reflect.DeepEqual(canonicalOld, canonicalNew)
}
