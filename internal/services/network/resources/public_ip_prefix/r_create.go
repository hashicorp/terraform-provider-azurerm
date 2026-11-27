// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package public_ip_prefix

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePublicIpPrefixCreate(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).PublicIPPrefixes
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := publicipprefixes.NewPublicIPPrefixID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, publicipprefixes.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_public_ip_prefix", id.ID())
		}
	}

	publicIpPrefix := publicipprefixes.PublicIPPrefix{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Sku: &publicipprefixes.PublicIPPrefixSku{
			Name: pointer.ToEnum[publicipprefixes.PublicIPPrefixSkuName](d.Get("sku").(string)),
			Tier: pointer.ToEnum[publicipprefixes.PublicIPPrefixSkuTier](d.Get("sku_tier").(string)),
		},
		Properties: &publicipprefixes.PublicIPPrefixPropertiesFormat{
			PrefixLength:           pointer.To(int64(d.Get("prefix_length").(int))),
			PublicIPAddressVersion: pointer.ToEnum[publicipprefixes.IPVersion](d.Get("ip_version").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if customIpPrefixId := d.Get("custom_ip_prefix_id").(string); customIpPrefixId != "" {
		publicIpPrefix.Properties.CustomIPPrefix = &publicipprefixes.SubResource{
			Id: pointer.To(customIpPrefixId),
		}
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		publicIpPrefix.Zones = &zones
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, publicIpPrefix, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourcePublicIpPrefixRead(d, meta)
}
