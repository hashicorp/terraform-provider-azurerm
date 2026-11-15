// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type FunctionAppFunctionKeysDataSource struct{}

func TestAccFunctionAppFunctionKeysDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app_function_keys", "test")
	d := FunctionAppFunctionKeysDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("default_key").IsNotEmpty(),
				check.That(data.ResourceName).Key("keys.%").HasValue("1"),
				check.That(data.ResourceName).Key("keys.default").IsNotEmpty(),
				acceptance.TestCheckResourceAttrPair(data.ResourceName, "default_key", data.ResourceName, "keys.default"),
			),
		},
	})
}

func TestAccFunctionAppFunctionKeysDataSource_functionNotFound(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app_function_keys", "test")
	d := FunctionAppFunctionKeysDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config:      d.functionNotFound(data),
			ExpectError: regexp.MustCompile("was not found"),
		},
	})
}

func (FunctionAppFunctionKeysDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_function_app_function_keys" "test" {
  function_app_id = azurerm_function_app_function.test.function_app_id
  function_name   = azurerm_function_app_function.test.name
}
`, FunctionAppFunctionResource{}.withLocalFiles(data))
}

func (FunctionAppFunctionKeysDataSource) functionNotFound(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_function_app_function_keys" "test" {
  function_app_id = azurerm_linux_function_app.test.id
  function_name   = "doesNotExist"
}
`, FunctionAppFunctionResource{}.templateLinux(data))
}
