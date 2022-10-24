package cdn_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontdoorSecretResourceDataSource struct {
	DoNotRunFrontDoorCustomDomainTests string
}

// NOTE: This is currently not testable due to the cert requirements of the service
func TestAccCdnFrontDoorSecretDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_secret", "test")
	r := CdnFrontdoorSecretResource{os.Getenv("ARM_TEST_DO_NOT_RUN_CDN_FRONT_DOOR_CUSTOM_DOMAIN")}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cdn_frontdoor_profile_id").MatchesOtherKey(check.That("azurerm_cdn_frontdoor_profile.test").Key("id")),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontdoorSecretResourceDataSource) preCheck(t *testing.T) {
	if r.DoNotRunFrontDoorCustomDomainTests == "" {
		t.Skipf("`ARM_TEST_DO_NOT_RUN_CDN_FRONT_DOOR_CUSTOM_DOMAIN` must be set for acceptance tests")
	}

	if strings.EqualFold(r.DoNotRunFrontDoorCustomDomainTests, "true") {
		t.Skipf("`data.azurerm_cdn_frontdoor_secret` currently is not testable due to service requirements")
	}
}

func (r CdnFrontdoorSecretResourceDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "accTestProfile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_secret" "test" {
  name                     = "accTestSecret-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  secret {
    customer_certificate {
      key_vault_certificate_id = azurerm_key_vault_certificate.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnFrontdoorSecretResourceDataSource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_cdn_frontdoor_secret" "test" {
  name                = azurerm_cdn_frontdoor_secret.test.name
  profile_name        = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
}
`, template)
}
