// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
)

type FrameworkServiceRegistration interface {
	Actions() []func() action.Action

	FrameworkResources() []FrameworkWrappedResource

	FrameworkDataSources() []FrameworkWrappedDataSource

	EphemeralResources() []func() ephemeral.EphemeralResource

	ListResources() []func() list.ListResource
}
