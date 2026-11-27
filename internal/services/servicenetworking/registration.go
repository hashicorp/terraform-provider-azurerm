// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicenetworking

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/resources/application_load_balancer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/resources/application_load_balancer_frontend"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/resources/application_load_balancer_security_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/resources/application_load_balancer_subnet_association"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/service-networking"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		application_load_balancer.Resource{},
		application_load_balancer_subnet_association.Resource{},
		application_load_balancer_frontend.Resource{},
		application_load_balancer_security_policy.Resource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Service Networking"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Service Networking",
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
