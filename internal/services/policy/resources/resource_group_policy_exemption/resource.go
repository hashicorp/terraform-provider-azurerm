// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group_policy_exemption

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceResourceGroupPolicyExemptionCreateUpdate,
		Read:   resourceResourceGroupPolicyExemptionRead,
		Update: resourceResourceGroupPolicyExemptionCreateUpdate,
		Delete: resourceResourceGroupPolicyExemptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ResourceGroupPolicyExemptionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceGroupPolicyExemptionSchema(),
	}
}
