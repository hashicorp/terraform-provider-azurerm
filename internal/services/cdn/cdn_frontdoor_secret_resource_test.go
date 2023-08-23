// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorSecretResource struct {
	DoNotRunFrontDoorCustomDomainTests string
}

func TestAccCdnFrontDoorSecret_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_secret", "test")
	r := CdnFrontdoorSecretResource{os.Getenv("ARM_TEST_DO_NOT_RUN_CDN_FRONT_DOOR_CUSTOM_DOMAIN")}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorSecret_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_secret", "test")
	r := CdnFrontdoorSecretResource{os.Getenv("ARM_TEST_DO_NOT_RUN_CDN_FRONT_DOOR_CUSTOM_DOMAIN")}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnFrontDoorSecret_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_secret", "test")
	r := CdnFrontdoorSecretResource{os.Getenv("ARM_TEST_DO_NOT_RUN_CDN_FRONT_DOOR_CUSTOM_DOMAIN")}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontdoorSecretResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorSecretID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorSecretsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecretName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CdnFrontdoorSecretResource) preCheck(t *testing.T) {
	if r.DoNotRunFrontDoorCustomDomainTests == "" {
		t.Skipf("`ARM_TEST_DO_NOT_RUN_CDN_FRONT_DOOR_CUSTOM_DOMAIN` must be set for acceptance tests")
	}

	if strings.EqualFold(r.DoNotRunFrontDoorCustomDomainTests, "true") {
		t.Skipf("`azurerm_cdn_frontdoor_secret` currently is not testable due to service requirements")
	}
}

func (r CdnFrontdoorSecretResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%d"
  location = "%s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "accTestProfile-%d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CdnFrontdoorSecretResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_secret" "test" {
  name                     = "accTestSecret-%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  secret {
    customer_certificate {
      key_vault_certificate_id = azurerm_key_vault_certificate.test.id
    }
  }
}
`, template, data.RandomInteger)
}

func (r CdnFrontdoorSecretResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_secret" "import" {
  name                     = azurerm_cdn_frontdoor_secret.test.name
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  secret {
    customer_certificate {
      key_vault_certificate_id = azurerm_key_vault_certificate.test.id
    }
  }
}
`, config)
}

func (r CdnFrontdoorSecretResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_cdn_frontdoor_secret" "test" {
  name                     = "accTestSecret-%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  secret {
    customer_certificate {
      key_vault_certificate_id = azurerm_key_vault_certificate.test.versionless_id
    }
  }
}
`, template, data.RandomInteger)
}
