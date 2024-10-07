// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/paloalto"
}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

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
		NextGenerationFirewallVNetLocalRulestackResource{},
		NextGenerationFirewallVNetPanoramaResource{},
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Palo Alto",
	}
}
