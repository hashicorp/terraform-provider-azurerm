// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_local_rulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_local_rulestack_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_local_rulestack_fqdn_list"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_local_rulestack_outbound_trust_certificate_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_local_rulestack_outbound_untrust_certificate_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_local_rulestack_prefix_list"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_local_rulestack_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_next_generation_firewall_virtual_hub_local_rulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_next_generation_firewall_virtual_hub_panorama"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_next_generation_firewall_virtual_network_local_rulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_next_generation_firewall_virtual_network_panorama"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/resources/palo_alto_virtual_network_appliance"
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
		palo_alto_local_rulestack.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		palo_alto_local_rulestack.Resource{},
		palo_alto_local_rulestack_certificate.Resource{},
		palo_alto_local_rulestack_fqdn_list.Resource{},
		palo_alto_local_rulestack_outbound_trust_certificate_association.Resource{},
		palo_alto_local_rulestack_outbound_untrust_certificate_association.Resource{},
		palo_alto_local_rulestack_prefix_list.Resource{},
		palo_alto_local_rulestack_rule.Resource{},
		palo_alto_virtual_network_appliance.Resource{},
		palo_alto_next_generation_firewall_virtual_hub_local_rulestack.Resource{},
		palo_alto_next_generation_firewall_virtual_hub_panorama.Resource{},
		palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager.Resource{},
		palo_alto_next_generation_firewall_virtual_network_local_rulestack.Resource{},
		palo_alto_next_generation_firewall_virtual_network_panorama.Resource{},
		palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager.Resource{},
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
