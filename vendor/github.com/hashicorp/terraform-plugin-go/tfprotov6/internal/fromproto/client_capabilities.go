// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func ValidateResourceConfigClientCapabilities(in *tfplugin6.ClientCapabilities) *tfprotov6.ValidateResourceConfigClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ValidateResourceConfigClientCapabilities{
		WriteOnlyAttributesAllowed: in.WriteOnlyAttributesAllowed,
	}

	return resp
}

func ConfigureProviderClientCapabilities(in *tfplugin6.ClientCapabilities) *tfprotov6.ConfigureProviderClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ConfigureProviderClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadDataSourceClientCapabilities(in *tfplugin6.ClientCapabilities) *tfprotov6.ReadDataSourceClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ReadDataSourceClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadResourceClientCapabilities(in *tfplugin6.ClientCapabilities) *tfprotov6.ReadResourceClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ReadResourceClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func PlanResourceChangeClientCapabilities(in *tfplugin6.ClientCapabilities) *tfprotov6.PlanResourceChangeClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.PlanResourceChangeClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ImportResourceStateClientCapabilities(in *tfplugin6.ClientCapabilities) *tfprotov6.ImportResourceStateClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ImportResourceStateClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func OpenEphemeralResourceClientCapabilities(in *tfplugin6.ClientCapabilities) *tfprotov6.OpenEphemeralResourceClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.OpenEphemeralResourceClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}
