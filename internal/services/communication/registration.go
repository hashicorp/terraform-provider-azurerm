// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/resources/communication_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/resources/communication_service_email_domain_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/resources/email_communication_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/resources/email_communication_service_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/resources/email_communication_service_domain_sender_username"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/communication"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Communication"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Communication",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		communication_service.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		communication_service.Resource{},
		email_communication_service_domain.Resource{},
		email_communication_service_domain_sender_username.Resource{},
		email_communication_service.Resource{},
		communication_service_email_domain_association.Resource{},
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
