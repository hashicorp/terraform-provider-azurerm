// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bastion_host

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/bastionhosts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBastionHostRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.BastionHostsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := bastionhosts.ParseBastionHostID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.BastionHostName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku", string(*sku.Name))
		}

		d.Set("zones", zones.FlattenUntyped(model.Zones))

		if props := model.Properties; props != nil {
			d.Set("dns_name", props.DnsName)
			d.Set("scale_units", props.ScaleUnits)
			d.Set("file_copy_enabled", props.EnableFileCopy)
			d.Set("ip_connect_enabled", props.EnableIPConnect)
			d.Set("kerberos_enabled", props.EnableKerberos)
			d.Set("shareable_link_enabled", props.EnableShareableLink)
			d.Set("tunneling_enabled", props.EnableTunneling)
			d.Set("session_recording_enabled", props.EnableSessionRecording)
			d.Set("private_only_enabled", props.EnablePrivateOnlyBastion)

			virtualNetworkId := ""
			if vnet := props.VirtualNetwork; vnet != nil {
				vnetId, err := commonids.ParseVirtualNetworkID(pointer.From(vnet.Id))
				if err != nil {
					return err
				}
				virtualNetworkId = vnetId.ID()
			}
			d.Set("virtual_network_id", virtualNetworkId)

			copyPasteEnabled := true
			if props.DisableCopyPaste != nil {
				copyPasteEnabled = !*props.DisableCopyPaste
			}
			d.Set("copy_paste_enabled", copyPasteEnabled)

			if ipConfigs := props.IPConfigurations; ipConfigs != nil {
				if err := d.Set("ip_configuration", flattenBastionHostIPConfiguration(ipConfigs)); err != nil {
					return fmt.Errorf("flattening `ip_configuration`: %+v", err)
				}
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenBastionHostIPConfiguration(ipConfigs *[]bastionhosts.BastionHostIPConfiguration) []any {
	result := make([]any, 0)
	if ipConfigs == nil {
		return result
	}

	for _, config := range *ipConfigs {
		ipConfig := make(map[string]any)

		if config.Name != nil {
			ipConfig["name"] = *config.Name
		}

		if props := config.Properties; props != nil {
			ipConfig["subnet_id"] = pointer.From(props.Subnet.Id)

			publicIpId := ""
			if pip := props.PublicIPAddress; pip != nil && pip.Id != nil {
				publicIpId = *pip.Id
			}
			ipConfig["public_ip_address_id"] = publicIpId
		}

		result = append(result, ipConfig)
	}
	return result
}

var skuWeight = map[string]int8{
	"Developer": 1,
	"Basic":     2,
	"Standard":  3,
	"Premium":   4,
}
