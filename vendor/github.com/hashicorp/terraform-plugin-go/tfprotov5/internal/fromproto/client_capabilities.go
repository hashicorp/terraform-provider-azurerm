// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func ConfigureProviderClientCapabilities(in *tfplugin5.ClientCapabilities) *tfprotov5.ConfigureProviderClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.ConfigureProviderClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadDataSourceClientCapabilities(in *tfplugin5.ClientCapabilities) *tfprotov5.ReadDataSourceClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.ReadDataSourceClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadResourceClientCapabilities(in *tfplugin5.ClientCapabilities) *tfprotov5.ReadResourceClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.ReadResourceClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func PlanResourceChangeClientCapabilities(in *tfplugin5.ClientCapabilities) *tfprotov5.PlanResourceChangeClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.PlanResourceChangeClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ImportResourceStateClientCapabilities(in *tfplugin5.ClientCapabilities) *tfprotov5.ImportResourceStateClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.ImportResourceStateClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}
