// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func DataSourceDeferred(fw *datasource.Deferred) *tfprotov5.Deferred {
	if fw == nil {
		return nil
	}
	return &tfprotov5.Deferred{
		Reason: tfprotov5.DeferredReason(fw.Reason),
	}
}

func ResourceDeferred(fw *resource.Deferred) *tfprotov5.Deferred {
	if fw == nil {
		return nil
	}
	return &tfprotov5.Deferred{
		Reason: tfprotov5.DeferredReason(fw.Reason),
	}
}

func EphemeralResourceDeferred(fw *ephemeral.Deferred) *tfprotov5.Deferred {
	if fw == nil {
		return nil
	}
	return &tfprotov5.Deferred{
		Reason: tfprotov5.DeferredReason(fw.Reason),
	}
}

func ActionDeferred(fw *action.Deferred) *tfprotov5.Deferred {
	if fw == nil {
		return nil
	}
	return &tfprotov5.Deferred{
		Reason: tfprotov5.DeferredReason(fw.Reason),
	}
}
