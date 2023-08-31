// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontDoorCustomDomainDataSource struct{}

func TestAccCdnFrontDoorCustomDomainDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_custom_domain", "test")
	d := CdnFrontDoorCustomDomainDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("dns_zone_id").Exists(),
				check.That(data.ResourceName).Key("host_name").Exists(),
				check.That(data.ResourceName).Key("cdn_frontdoor_profile_id").Exists(),
				check.That(data.ResourceName).Key("tls.0.cdn_frontdoor_secret_id").IsEmpty(),
				check.That(data.ResourceName).Key("tls.0.certificate_type").Exists(),
				check.That(data.ResourceName).Key("tls.0.minimum_tls_version").Exists(),
				check.That(data.ResourceName).Key("expiration_date").Exists(),
				check.That(data.ResourceName).Key("validation_token").Exists(),
			),
		},
	})
}

func (CdnFrontDoorCustomDomainDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                = azurerm_cdn_frontdoor_custom_domain.test.name
  profile_name        = azurerm_cdn_frontdoor_profile.test.name
  resource_group_name = azurerm_cdn_frontdoor_profile.test.resource_group_name
}
`, CdnFrontDoorCustomDomainResource{}.complete(data))
}
