// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nat_gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var natGatewayResourceName = "azurerm_nat_gateway"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNatGatewayCreate,
		Read:   resourceNatGatewayRead,
		Update: resourceNatGatewayUpdate,
		Delete: resourceNatGatewayDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&natgateways.NatGatewayId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&natgateways.NatGatewayId{}),
		},

		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, _ any) error {
			if diff.Get("sku_name").(string) == string(natgateways.NatGatewaySkuNameStandardVTwo) {
				if !diff.GetRawConfig().AsValueMap()["zones"].IsNull() {
					return fmt.Errorf("%s resources with `sku_name` set to `%s` are zone-redundant by default, Azure automatically deploys across all available zones. The `zones` argument must be omitted", natGatewayResourceName, natgateways.NatGatewaySkuNameStandardVTwo)
				}
			}

			return nil
		},

		Schema: resourceNatGatewaySchema(),
	}
}
