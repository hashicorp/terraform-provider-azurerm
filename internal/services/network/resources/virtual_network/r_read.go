// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network

import (
	"bytes"
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if err := resourceVirtualNetworkFlatten(d, *id, resp.Model); err != nil {
		return fmt.Errorf("encoding %s: %+v", *id, err)
	}

	return nil
}

func resourceAzureSubnetHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		buf.WriteString(m["name"].(string))
		if v, ok := m["address_prefix"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["security_group"]; ok {
			buf.WriteString(v.(string))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func getExistingSubnet(ctx context.Context, client virtualnetworks.VirtualNetworksClient, id commonids.VirtualNetworkId, subnetName string) (*virtualnetworks.Subnet, error) {
	resp, err := client.Get(ctx, id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		// The Subnet doesn't exist when the Virtual Network doesn't exist
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(virtualnetworks.Subnet{}), nil
		}
		// raise an error if there was an issue other than 404 in getting subnet properties
		return nil, err
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			for _, subnet := range pointer.From(props.Subnets) {
				if subnetName == pointer.From(subnet.Name) {
					// Return it directly rather than copy the fields to prevent potential uncovered properties (for example, `ServiceEndpoints` mentioned in #1619)
					return pointer.To(subnet), nil
				}
			}
		}
	}

	return pointer.To(virtualnetworks.Subnet{}), nil
}
