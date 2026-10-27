// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_zone_virtual_network_link

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateDnsZoneVirtualNetworkLinkCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).PrivateDns.VirtualNetworkLinksClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	privateDNSZoneID, err := privatezones.ParsePrivateDnsZoneID(d.Get("private_dns_zone_id").(string))
	if err != nil {
		return err
	}

	id := virtualnetworklinks.NewVirtualNetworkLinkID(privateDNSZoneID.SubscriptionId, privateDNSZoneID.ResourceGroupName, privateDNSZoneID.PrivateDnsZoneName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_private_dns_zone_virtual_network_link", id.ID())
			}
		}
	}

	parameters := virtualnetworklinks.VirtualNetworkLink{
		Location: pointer.To("global"),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Properties: &virtualnetworklinks.VirtualNetworkLinkProperties{
			VirtualNetwork: &virtualnetworklinks.SubResource{
				Id: pointer.To(d.Get("virtual_network_id").(string)),
			},
			RegistrationEnabled: pointer.To(d.Get("registration_enabled").(bool)),
		},
	}

	if v, ok := d.GetOk("resolution_policy"); ok {
		parameters.Properties.ResolutionPolicy = pointer.ToEnum[virtualnetworklinks.ResolutionPolicy](v.(string))
	}

	options := virtualnetworklinks.CreateOrUpdateOperationOptions{
		IfMatch:     pointer.To(""),
		IfNoneMatch: pointer.To(""),
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, options, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, options); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourcePrivateDnsZoneVirtualNetworkLinkRead(d, meta)
}
