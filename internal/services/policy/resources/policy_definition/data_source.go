// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_definition

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: helpers.PolicyDefinitionReadFunc(false),

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: helpers.PolicyDefinitionDataSourceSchema(),
	}
}
