package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type AzureFrontDoorCachePurgeAction struct{}

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

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%[3]s", azurerm_dns_zone.test.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
