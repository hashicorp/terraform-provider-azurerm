// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistration                 = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/paloalto"
}

func (r Registration) Name() string {
	return "Palo Alto"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		LocalRulestackDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		LocalRuleStack{},
		LocalRuleStackCertificate{},
		LocalRulestackFQDNList{},
		LocalRulestackOutboundTrustCertificateAssociationResource{},
		LocalRulestackOutboundUnTrustCertificateAssociationResource{},
		LocalRuleStackPrefixList{},
		LocalRuleStackRule{},
		NetworkVirtualApplianceResource{},
		NextGenerationFirewallVHubLocalRuleStackResource{},
		NextGenerationFirewallVHubPanoramaResource{},
		NextGenerationFirewallVHubStrataCloudManagerResource{},
		NextGenerationFirewallVNetLocalRulestackResource{},
		NextGenerationFirewallVNetPanoramaResource{},
		NextGenerationFirewallVNetStrataCloudManagerResource{},
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Palo Alto",
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
