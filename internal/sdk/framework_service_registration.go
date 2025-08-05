// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type FrameworkTypedServiceRegistration interface {
	TypedServiceRegistration

	FrameworkResources() []func() resource.Resource

	FrameworkDataSources() []func() datasource.DataSource

	EphemeralResources() []func() ephemeral.EphemeralResource

	ListResources() []func() list.ListResource
}
