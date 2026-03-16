// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package extendedlocation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ExtendedLocationCustomLocationDataSource struct{}

func TestAccExtendedLocationCustomLocationDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_extended_location_custom_location", "test")
	d := ExtendedLocationCustomLocationDataSource{}

	credential, privateKey, publicKey := ExtendedLocationCustomLocationResource{}.getCredentials(t)

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.basic(data, credential, privateKey, publicKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").IsNotEmpty(),
				check.That(data.ResourceName).Key("host_type").HasValue("Kubernetes"),
				check.That(data.ResourceName).Key("cluster_extension_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("display_name").IsNotEmpty(),
				check.That(data.ResourceName).Key("namespace").IsNotEmpty(),
				check.That(data.ResourceName).Key("host_resource_id").IsNotEmpty(),
			),
		},
	})
}

func (d ExtendedLocationCustomLocationDataSource) basic(data acceptance.TestData, credential string, privateKey string, publicKey string) string {
	return fmt.Sprintf(`
%s

data "azurerm_extended_location_custom_location" "test" {
  name                = azurerm_extended_location_custom_location.test.name
  resource_group_name = azurerm_extended_location_custom_location.test.resource_group_name
}
`, ExtendedLocationCustomLocationResource{}.complete(data, credential, privateKey, publicKey))
}
