// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func DataSourceDeferred(fw *datasource.Deferred) *tfprotov6.Deferred {
	if fw == nil {
		return nil
	}
	return &tfprotov6.Deferred{
		Reason: tfprotov6.DeferredReason(fw.Reason),
	}
}

func ResourceDeferred(fw *resource.Deferred) *tfprotov6.Deferred {
	if fw == nil {
		return nil
	}
	return &tfprotov6.Deferred{
		Reason: tfprotov6.DeferredReason(fw.Reason),
	}
}

func EphemeralResourceDeferred(fw *ephemeral.Deferred) *tfprotov6.Deferred {
	if fw == nil {
		return nil
	}
	return &tfprotov6.Deferred{
		Reason: tfprotov6.DeferredReason(fw.Reason),
	}
}

func ActionDeferred(fw *action.Deferred) *tfprotov6.Deferred {
	if fw == nil {
		return nil
	}
	return &tfprotov6.Deferred{
		Reason: tfprotov6.DeferredReason(fw.Reason),
	}
}
