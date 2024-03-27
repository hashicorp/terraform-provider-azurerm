// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StaticWebAppDataSource struct{}

func TestAccAzureStaticWebAppDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_static_web_app", "test")
	r := StaticWebAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("default_host_name").Exists(),
				check.That(data.ResourceName).Key("api_key").Exists(),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acceptance"),
				check.That(data.ResourceName).Key("basic_auth.#").HasValue("1"),
				check.That(data.ResourceName).Key("basic_auth.0.environments").HasValue("AllEnvironments"),
				check.That(data.ResourceName).Key("identity.#").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
			),
		},
	})
}

func (StaticWebAppDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

data "azurerm_static_web_app" "test" {
  name                = azurerm_static_web_app.test.name
  resource_group_name = azurerm_static_web_app.test.resource_group_name
}
`, StaticWebAppResource{}.complete(data))
}
