// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestAccPaloAltoNextGenerationFirewallVHubLocalRulestackDataSource_basic(t *testing.T) {
	resourceType := "azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack"
	data := acceptance.BuildTestData(t, "data."+resourceType, "test")

	data.DataSourceTest(t, []acceptance.TestStep{{
		Config: nextGenerationFirewallDataSourceConfig(NextGenerationFirewallVWanResource{}.complete(data), resourceType),
		Check: acceptance.ComposeTestCheckFunc(
			check.That(data.ResourceName).Key("destination_nat.0.name").HasValue("testDNAT-1"),
			check.That(data.ResourceName).Key("marketplace_offer_id").HasValue("pan_swfw_cloud_ngfw"),
			check.That(data.ResourceName).Key("network_profile.0.virtual_hub_id").Exists(),
			check.That(data.ResourceName).Key("plan_id").HasValue("panw-cngfw-payg"),
			check.That(data.ResourceName).Key("rulestack_id").Exists(),
		),
	}})
}

func TestAccPaloAltoNextGenerationFirewallVNetLocalRulestackDataSource_basic(t *testing.T) {
	resourceType := "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack"
	data := acceptance.BuildTestData(t, "data."+resourceType, "test")

	data.DataSourceTest(t, []acceptance.TestStep{{
		Config: nextGenerationFirewallDataSourceConfig(NextGenerationFirewallVnetResource{}.complete(data), resourceType),
		Check: acceptance.ComposeTestCheckFunc(
			check.That(data.ResourceName).Key("destination_nat.0.name").HasValue("testDNAT-1"),
			check.That(data.ResourceName).Key("marketplace_offer_id").HasValue("pan_swfw_cloud_ngfw"),
			check.That(data.ResourceName).Key("network_profile.0.vnet_configuration.0.virtual_network_id").Exists(),
			check.That(data.ResourceName).Key("plan_id").HasValue("panw-cngfw-payg"),
			check.That(data.ResourceName).Key("rulestack_id").Exists(),
		),
	}})
}

func TestAccPaloAltoNextGenerationFirewallVHubPanoramaDataSource_basic(t *testing.T) {
	if os.Getenv("ARM_PALO_ALTO_PANORAMA_CONFIG") == "" {
		t.Skip("skipping as Palo Alto Panorama config not set in `ARM_PALO_ALTO_PANORAMA_CONFIG`")
	}

	resourceType := "azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama"
	data := acceptance.BuildTestData(t, "data."+resourceType, "test")

	data.DataSourceTest(t, []acceptance.TestStep{{
		Config: nextGenerationFirewallDataSourceConfig(NextGenerationFirewallVHubPanoramaResource{}.complete(data), resourceType),
		Check: acceptance.ComposeTestCheckFunc(
			check.That(data.ResourceName).Key("location").Exists(),
			check.That(data.ResourceName).Key("marketplace_offer_id").HasValue("pan_swfw_cloud_ngfw"),
			check.That(data.ResourceName).Key("network_profile.0.virtual_hub_id").Exists(),
			check.That(data.ResourceName).Key("panorama_base64_config").Exists(),
			check.That(data.ResourceName).Key("plan_id").HasValue("panw-cngfw-payg"),
		),
	}})
}

func TestAccPaloAltoNextGenerationFirewallVNetPanoramaDataSource_basic(t *testing.T) {
	if os.Getenv("ARM_PALO_ALTO_PANORAMA_CONFIG") == "" {
		t.Skip("skipping as Palo Alto Panorama config not set in `ARM_PALO_ALTO_PANORAMA_CONFIG`")
	}

	resourceType := "azurerm_palo_alto_next_generation_firewall_virtual_network_panorama"
	data := acceptance.BuildTestData(t, "data."+resourceType, "test")

	data.DataSourceTest(t, []acceptance.TestStep{{
		Config: nextGenerationFirewallDataSourceConfig(NextGenerationFirewallVNetPanoramaResource{}.complete(data), resourceType),
		Check: acceptance.ComposeTestCheckFunc(
			check.That(data.ResourceName).Key("location").Exists(),
			check.That(data.ResourceName).Key("marketplace_offer_id").HasValue("pan_swfw_cloud_ngfw"),
			check.That(data.ResourceName).Key("network_profile.0.vnet_configuration.0.virtual_network_id").Exists(),
			check.That(data.ResourceName).Key("panorama_base64_config").Exists(),
			check.That(data.ResourceName).Key("plan_id").HasValue("panw-cngfw-payg"),
		),
	}})
}

func TestAccPaloAltoNextGenerationFirewallVHubStrataCloudManagerDataSource_basic(t *testing.T) {
	if os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME") == "" {
		t.Skip("skipping as Palo Alto Strata Cloud Manager tenant name not set in `ARM_PALO_ALTO_SCM_TENANT_NAME`")
	}

	resourceType := "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager"
	data := acceptance.BuildTestData(t, "data."+resourceType, "test")

	data.DataSourceTest(t, []acceptance.TestStep{{
		Config: nextGenerationFirewallDataSourceConfig(NextGenerationFirewallVHubStrataCloudManagerResource{}.complete(data), resourceType),
		Check: acceptance.ComposeTestCheckFunc(
			check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
			check.That(data.ResourceName).Key("location").Exists(),
			check.That(data.ResourceName).Key("network_profile.0.virtual_hub_id").Exists(),
			check.That(data.ResourceName).Key("plan_id").HasValue("panw-cngfw-payg"),
			check.That(data.ResourceName).Key("strata_cloud_manager_tenant_name").HasValue(os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME")),
		),
	}})
}

func TestAccPaloAltoNextGenerationFirewallVNetStrataCloudManagerDataSource_basic(t *testing.T) {
	if os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME") == "" {
		t.Skip("skipping as Palo Alto Strata Cloud Manager tenant name not set in `ARM_PALO_ALTO_SCM_TENANT_NAME`")
	}

	resourceType := "azurerm_palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager"
	data := acceptance.BuildTestData(t, "data."+resourceType, "test")

	data.DataSourceTest(t, []acceptance.TestStep{{
		Config: nextGenerationFirewallDataSourceConfig(NextGenerationFirewallVNetStrataCloudManagerResource{}.complete(data), resourceType),
		Check: acceptance.ComposeTestCheckFunc(
			check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
			check.That(data.ResourceName).Key("location").Exists(),
			check.That(data.ResourceName).Key("network_profile.0.vnet_configuration.0.virtual_network_id").Exists(),
			check.That(data.ResourceName).Key("plan_id").HasValue("panw-cngfw-payg"),
			check.That(data.ResourceName).Key("strata_cloud_manager_tenant_name").HasValue(os.Getenv("ARM_PALO_ALTO_SCM_TENANT_NAME")),
		),
	}})
}

func TestPaloAltoNextGenerationFirewallDataSourceSchemasMatchResources(t *testing.T) {
	tests := []struct {
		name       string
		resource   sdk.Resource
		dataSource sdk.DataSource
	}{
		{"virtual hub local rulestack", paloalto.NextGenerationFirewallVHubLocalRuleStackResource{}, paloalto.NextGenerationFirewallVHubLocalRulestackDataSource{}},
		{"virtual hub panorama", paloalto.NextGenerationFirewallVHubPanoramaResource{}, paloalto.NextGenerationFirewallVHubPanoramaDataSource{}},
		{"virtual hub strata cloud manager", paloalto.NextGenerationFirewallVHubStrataCloudManagerResource{}, paloalto.NextGenerationFirewallVHubStrataCloudManagerDataSource{}},
		{"virtual network local rulestack", paloalto.NextGenerationFirewallVNetLocalRulestackResource{}, paloalto.NextGenerationFirewallVNetLocalRulestackDataSource{}},
		{"virtual network panorama", paloalto.NextGenerationFirewallVNetPanoramaResource{}, paloalto.NextGenerationFirewallVNetPanoramaDataSource{}},
		{"virtual network strata cloud manager", paloalto.NextGenerationFirewallVNetStrataCloudManagerResource{}, paloalto.NextGenerationFirewallVNetStrataCloudManagerDataSource{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.resource.ResourceType() != test.dataSource.ResourceType() {
				t.Fatalf("resource type mismatch: resource %q, data source %q", test.resource.ResourceType(), test.dataSource.ResourceType())
			}

			resourceSchema := mergeSchemas(test.resource.Arguments(), test.resource.Attributes())
			dataSourceSchema := mergeSchemas(test.dataSource.Arguments(), test.dataSource.Attributes())
			assertSchemaShapeMatches(t, test.resource.ResourceType(), resourceSchema, dataSourceSchema)
		})
	}
}

func nextGenerationFirewallDataSourceConfig(resourceConfig, resourceType string) string {
	return fmt.Sprintf(`
%[1]s

data %[2]q "test" {
  name                = %[2]s.test.name
  resource_group_name = %[2]s.test.resource_group_name
}
`, resourceConfig, resourceType)
}

func mergeSchemas(first, second map[string]*pluginsdk.Schema) map[string]*pluginsdk.Schema {
	result := make(map[string]*pluginsdk.Schema, len(first)+len(second))
	for key, value := range first {
		result[key] = value
	}
	for key, value := range second {
		result[key] = value
	}
	return result
}

func assertSchemaShapeMatches(t *testing.T, path string, resourceSchema, dataSourceSchema map[string]*pluginsdk.Schema) {
	t.Helper()

	for name, resourceField := range resourceSchema {
		dataSourceField, ok := dataSourceSchema[name]
		if !ok {
			t.Errorf("%s: data source does not expose resource field %q", path, name)
			continue
		}

		if resourceField.Type != dataSourceField.Type {
			t.Errorf("%s.%s: schema type mismatch: resource %v, data source %v", path, name, resourceField.Type, dataSourceField.Type)
		}

		resourceNested, resourceHasNested := resourceField.Elem.(*pluginsdk.Resource)
		dataSourceNested, dataSourceHasNested := dataSourceField.Elem.(*pluginsdk.Resource)
		if resourceHasNested != dataSourceHasNested {
			t.Errorf("%s.%s: nested resource shape mismatch", path, name)
			continue
		}
		if resourceHasNested {
			assertSchemaShapeMatches(t, path+"."+name, resourceNested.Schema, dataSourceNested.Schema)
		}
	}

	for name := range dataSourceSchema {
		if _, ok := resourceSchema[name]; !ok {
			t.Errorf("%s: data source exposes field %q which is not exposed by the resource", path, name)
		}
	}
}
