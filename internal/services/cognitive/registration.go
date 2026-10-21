// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account_connection_account_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account_connection_account_managed_identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account_connection_api_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account_connection_custom_keys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account_connection_entra_id"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account_customer_managed_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account_project"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account_rai_blocklist"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_account_rai_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/resources/cognitive_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/cognitive-services"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Cognitive Services"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Cognitive Services",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_cognitive_account": cognitive_account.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_cognitive_account":                      cognitive_account.RegisterResource(),
		"azurerm_cognitive_account_customer_managed_key": cognitive_account_customer_managed_key.RegisterResource(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		cognitive_account_project.DataSource{},
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		cognitive_account_connection_account_key.Resource{},
		cognitive_account_connection_account_managed_identity.Resource{},
		cognitive_account_connection_api_key.Resource{},
		cognitive_account_connection_custom_keys.Resource{},
		cognitive_account_connection_entra_id.Resource{},
		cognitive_account_project.Resource{},
		cognitive_account_rai_policy.Resource{},
		cognitive_deployment.Resource{},
		cognitive_account_rai_blocklist.Resource{},
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
	return []sdk.FrameworkListWrappedResource{
		cognitive_account_connection_account_key.ListResource{},
		cognitive_account_connection_account_managed_identity.ListResource{},
		cognitive_account_connection_api_key.ListResource{},
		cognitive_account_connection_custom_keys.ListResource{},
		cognitive_account_connection_entra_id.ListResource{},
		cognitive_account.ListResource{},
		cognitive_deployment.ListResource{},
	}
}
