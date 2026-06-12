// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func ConfigureProviderClientCapabilities(in *tfprotov6.ConfigureProviderClientCapabilities) provider.ConfigureProviderClientCapabilities {
	if in == nil {
		// Client did not indicate any supported capabilities
		return provider.ConfigureProviderClientCapabilities{
			DeferralAllowed: false,
		}
	}

	return provider.ConfigureProviderClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}
}

func ReadDataSourceClientCapabilities(in *tfprotov6.ReadDataSourceClientCapabilities) datasource.ReadClientCapabilities {
	if in == nil {
		// Client did not indicate any supported capabilities
		return datasource.ReadClientCapabilities{
			DeferralAllowed: false,
		}
	}

	return datasource.ReadClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}
}

func ReadResourceClientCapabilities(in *tfprotov6.ReadResourceClientCapabilities) resource.ReadClientCapabilities {
	if in == nil {
		// Client did not indicate any supported capabilities
		return resource.ReadClientCapabilities{
			DeferralAllowed: false,
		}
	}

	return resource.ReadClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}
}

func ModifyPlanClientCapabilities(in *tfprotov6.PlanResourceChangeClientCapabilities) resource.ModifyPlanClientCapabilities {
	if in == nil {
		// Client did not indicate any supported capabilities
		return resource.ModifyPlanClientCapabilities{
			DeferralAllowed: false,
		}
	}

	return resource.ModifyPlanClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}
}

func ImportStateClientCapabilities(in *tfprotov6.ImportResourceStateClientCapabilities) resource.ImportStateClientCapabilities {
	if in == nil {
		// Client did not indicate any supported capabilities
		return resource.ImportStateClientCapabilities{
			DeferralAllowed: false,
		}
	}

	return resource.ImportStateClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}
}

func OpenEphemeralResourceClientCapabilities(in *tfprotov6.OpenEphemeralResourceClientCapabilities) ephemeral.OpenClientCapabilities {
	if in == nil {
		// Client did not indicate any supported capabilities
		return ephemeral.OpenClientCapabilities{
			DeferralAllowed: false,
		}
	}

	return ephemeral.OpenClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}
}

func ValidateResourceConfigClientCapabilities(in *tfprotov6.ValidateResourceConfigClientCapabilities) resource.ValidateConfigClientCapabilities {
	if in == nil {
		// Client did not indicate any supported capabilities
		return resource.ValidateConfigClientCapabilities{
			WriteOnlyAttributesAllowed: false,
		}
	}

	return resource.ValidateConfigClientCapabilities{
		WriteOnlyAttributesAllowed: in.WriteOnlyAttributesAllowed,
	}
}

func ModifyPlanActionClientCapabilities(in *tfprotov6.PlanActionClientCapabilities) action.ModifyPlanClientCapabilities {
	if in == nil {
		// Client did not indicate any supported capabilities
		return action.ModifyPlanClientCapabilities{
			DeferralAllowed: false,
		}
	}

	return action.ModifyPlanClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}
}
