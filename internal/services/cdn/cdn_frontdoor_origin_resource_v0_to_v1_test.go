// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccCdnFrontDoorOrigin_V0ToV1_4810 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v4.81.0 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccCdnFrontDoorOrigin_V0ToV1_4810(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_origin", "test")
	r := CdnFrontdoorOriginResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestrg-cdn-afdx-%[2]d/providers/Microsoft.Cdn/profiles/acctest-cdnfdprofile-%[2]d/originGroups/acctest-cdnfd-group-%[2]d/origins/acctest-cdnfdorigin-%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicV0(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestrg-cdn-afdx-%[2]d/providers/microsoft.cdn/profiles/acctest-cdnfdprofile-%[2]d/originGroups/acctest-cdnfd-group-%[2]d/origins/acctest-cdnfdorigin-%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestrg-cdn-afdx-%[2]d/providers/Microsoft.Cdn/profiles/acctest-cdnfdprofile-%[2]d/originGroups/acctest-cdnfd-group-%[2]d/origins/acctest-cdnfdorigin-%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "4.81.0")
}

func (r CdnFrontdoorOriginResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_origin" "test-import" {
  name                          = "acctest-cdnfdorigin-%[2]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}

removed {
  from = azurerm_cdn_frontdoor_origin.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_cdn_frontdoor_origin.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acctestrg-cdn-afdx-%[2]d/providers/microsoft.cdn/profiles/acctest-cdnfdprofile-%[2]d/originGroups/acctest-cdnfd-group-%[2]d/origins/acctest-cdnfdorigin-%[2]d"
}
`, r.template(data, "Standard_AzureFrontDoor", false), data.RandomInteger, data.Subscriptions.Primary)
}

func (r CdnFrontdoorOriginResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_cdn_frontdoor_origin" "test-import" {
  name                          = "acctest-cdnfdorigin-%[2]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = "www.contoso.com"
  priority                       = 1
  weight                         = 1
}
`, r.template(data, "Standard_AzureFrontDoor", false), data.RandomInteger, data.Subscriptions.Primary)
}
