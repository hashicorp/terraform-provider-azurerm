// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_security_partner_provider

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/securitypartnerproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubSecurityPartnerProviderCreate,
		Read:   resourceVirtualHubSecurityPartnerProviderRead,
		Update: resourceVirtualHubSecurityPartnerProviderUpdate,
		Delete: resourceVirtualHubSecurityPartnerProviderDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := securitypartnerproviders.ParseSecurityPartnerProviderID(id)
			return err
		}),

		Schema: virtualHubSecurityPartnerProviderSchema(),
	}
}
