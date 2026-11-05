// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_security_policy

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorSecurityPolicyCreate,
		Read:   resourceCdnFrontdoorSecurityPolicyRead,
		Update: resourceCdnFrontdoorSecurityPolicyUpdate,
		Delete: resourceCdnFrontdoorSecurityPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := securitypolicies.ParseSecurityPolicyID(id)
			return err
		}),

		Schema: cdnFrontdoorSecurityPolicySchema(),
	}
}
