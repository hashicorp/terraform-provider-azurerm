// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_definition_built_in

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// dataSourcePolicyDefinitionBuiltIn read built-in policy definition only
func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about a Built-In Policy Definition. Retrieves Policy Definitions from your current subscription by default.",
		Read:        helpers.PolicyDefinitionReadFunc(true),

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: helpers.PolicyDefinitionDataSourceSchema(),
	}
}
