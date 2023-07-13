// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attestation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AttestationProviderDataSource struct{}

func TestAccAttestationProviderDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_attestation_provider", "test")
	resource := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AttestationProviderDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").MatchesOtherKey(check.That(resource.ResourceName).Key("name")),
				check.That(data.ResourceName).Key("resource_group_name").MatchesOtherKey(check.That(resource.ResourceName).Key("resource_group_name")),
				check.That(data.ResourceName).Key("location").MatchesOtherKey(check.That(resource.ResourceName).Key("location")),
				check.That(data.ResourceName).Key("attestation_uri").Exists(),
			),
		},
	})
}

func (AttestationProviderDataSource) basic(data acceptance.TestData) string {
	config := AttestationProviderResource{
		name: fmt.Sprintf("acctestap%s", data.RandomStringOfLength(10)),
	}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_attestation_provider" "test" {
  name                = azurerm_attestation_provider.test.name
  resource_group_name = azurerm_attestation_provider.test.resource_group_name
}
`, config)
}
