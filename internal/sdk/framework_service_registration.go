// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
)

type FrameworkServiceRegistration interface {
	Actions() []WrappedAction

	FrameworkResources() []FrameworkWrappedResource

	FrameworkDataSources() []FrameworkWrappedDataSource

	EphemeralResources() []func() ephemeral.EphemeralResource

	ListResources() []FrameworkListWrappedResource
}
