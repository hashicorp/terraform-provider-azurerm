// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_zone

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

var privateDnsZoneResourceName = "azurerm_private_dns_zone"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourcePrivateDnsZoneCreateUpdate,
		Read:     resourcePrivateDnsZoneRead,
		Update:   resourcePrivateDnsZoneCreateUpdate,
		Delete:   resourcePrivateDnsZoneDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&privatezones.PrivateDnsZoneId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&privatezones.PrivateDnsZoneId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: privateDnsZoneSchema(),
	}
}
