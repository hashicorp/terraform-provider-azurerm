// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package public_ip_prefix

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourcePublicIpPrefixCreate,
		Read:     resourcePublicIpPrefixRead,
		Update:   resourcePublicIpPrefixUpdate,
		Delete:   resourcePublicIpPrefixDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&publicipprefixes.PublicIPPrefixId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&publicipprefixes.PublicIPPrefixId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: publicIpPrefixSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
				skuTier := d.Get("sku_tier").(string)
				sku := d.Get("sku").(string)
				if strings.EqualFold(skuTier, string(publicipprefixes.PublicIPPrefixSkuTierGlobal)) && !strings.EqualFold(sku, string(publicipprefixes.PublicIPPrefixSkuNameStandard)) {
					return errors.New("`sku` must be set to `Standard` when `sku_tier` is set to `Global`")
				}
				return nil
			}),
		),
	}
}
