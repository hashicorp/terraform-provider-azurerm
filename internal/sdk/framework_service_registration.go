// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
)

type FrameworkServiceRegistration interface {
	FrameworkResources() []FrameworkWrappedResource

	FrameworkDataSources() []FrameworkWrappedDataSource

	EphemeralResources() []func() ephemeral.EphemeralResource
}
