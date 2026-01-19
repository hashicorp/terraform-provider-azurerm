package cdn_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type AzureFrontDoorCachePurgeAction struct{}

func (a AzureFrontDoorCachePurgeAction) preCheck(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME") == "" {
		t.Skipf("`ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME` must be set for acceptance tests!")
	}
	if os.Getenv("ARM_TEST_DNS_ZONE_NAME") == "" {
		t.Skipf("`ARM_TEST_DNS_ZONE_NAME` must be set for acceptance tests!")
	}
}

func TestAccAzureFrontDoorCachePurgeAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_front_door_cache_purge", "test")
	a := AzureFrontDoorCachePurgeAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.pathsOnly(data),
				Check:  nil, // TODO - terraform-plugin-testing release?
			},
		},
	})
}

func TestAccAzureFrontDoorCachePurgeAction_complete(t *testing.T) {
	t.Skip("skipping test: custom domains take a long, long time to deploy to an endpoint")
	data := acceptance.BuildTestData(t, "azurerm_cdn_front_door_cache_purge", "test")
	a := AzureFrontDoorCachePurgeAction{}
	a.preCheck(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.complete(data),
				Check:  nil, // TODO - terraform-plugin-testing release?
			},
		},
	})
}

func (a *AzureFrontDoorCachePurgeAction) pathsOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctest-cdnfdprofile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "acctest-cdnfdendpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "terraform_data" "trigger" {
  input = "trigger"
  lifecycle {
    action_trigger {
      events  = [before_create, before_update]
      actions = [action.azurerm_cdn_front_door_cache_purge.test]
    }
  }
}

action "azurerm_cdn_front_door_cache_purge" "test" {
  config {
    front_door_endpoint_id = azurerm_cdn_frontdoor_endpoint.test.id

    content_paths = [
      "/*"
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (a *AzureFrontDoorCachePurgeAction) complete(data acceptance.TestData) string {
	dnsZoneName := os.Getenv("ARM_TEST_DNS_ZONE_NAME")
	dnsZoneRG := os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctest-cdnfdprofile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "acctest-cdnfdendpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

data "azurerm_dns_zone" "test" {
  name                = "%[4]s"
  resource_group_name = "%[5]s"
}

locals {
  # Create a delegated child zone inside the test RG.
  # NOTE: ARM_TEST_DNS_ZONE_NAME / ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME must refer to a real, delegated parent zone.
  child_zone_label = "acctest%[1]d"
  child_zone_name  = join(".", [local.child_zone_label, data.azurerm_dns_zone.test.name])
}

resource "azurerm_dns_zone" "child" {
  name                = local.child_zone_name
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_ns_record" "delegation" {
  name                = local.child_zone_label
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  zone_name           = data.azurerm_dns_zone.test.name
  ttl                 = 300

  records = azurerm_dns_zone.child.name_servers
}

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  depends_on = [azurerm_dns_ns_record.delegation]

  dns_zone_id              = azurerm_dns_zone.child.id
  host_name                = join(".", ["%[3]s", azurerm_dns_zone.child.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

resource "azurerm_dns_txt_record" "validation" {
  depends_on = [azurerm_dns_ns_record.delegation]

  name                = join(".", ["_dnsauth", split(".", azurerm_cdn_frontdoor_custom_domain.test.host_name)[0]])
  zone_name           = azurerm_dns_zone.child.name
  resource_group_name = azurerm_resource_group.test.name
  ttl                 = 300

  record {
    value = azurerm_cdn_frontdoor_custom_domain.test.validation_token
  }
}

resource "terraform_data" "trigger" {
  input = "trigger"
  lifecycle {
    action_trigger {
      events  = [before_create, before_update]
      actions = [action.azurerm_cdn_front_door_cache_purge.test]
    }
  }
}

action "azurerm_cdn_front_door_cache_purge" "test" {
  config {
    front_door_endpoint_id = azurerm_cdn_frontdoor_endpoint.test.id

    content_paths = [
      "/*"
    ]

    domains = [
      azurerm_cdn_frontdoor_custom_domain.test.host_name
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, dnsZoneName, dnsZoneRG)
}
