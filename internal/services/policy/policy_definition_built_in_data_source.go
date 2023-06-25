package policy

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// dataSourceArmPolicyDefinitionBuiltIn read built-in policy definition only
func dataSourceArmPolicyDefinitionBuiltIn() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: policyDefinitionReadFunc(true),

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: policyDefinitionDataSourceSchema(),
	}
}
