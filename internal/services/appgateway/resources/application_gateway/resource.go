// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_gateway

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceApplicationGatewayCreate,
		Read:     resourceApplicationGatewayRead,
		Update:   resourceApplicationGatewayUpdate,
		Delete:   resourceApplicationGatewayDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&applicationgateways.ApplicationGatewayId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&applicationgateways.ApplicationGatewayId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: applicationGatewaySchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(applicationGatewayCustomizeDiff),
	}
}
