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

func (d FunctionAppFunctionKeysDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_function_app_function_keys" "test" {
  function_app_id = azurerm_function_app_function.test.function_app_id
  function_name   = azurerm_function_app_function.test.name
}
`, d.functionTemplate(data))
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

// functionTemplate provisions a Function with deployed code, since the Functions host only generates keys once it
// has successfully loaded the Function. `test_data` is intentionally omitted as it is not read back consistently.
func (FunctionAppFunctionKeysDataSource) functionTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_function" "test" {
  name            = "testAcc-FnAppFn-%[2]d"
  function_app_id = azurerm_windows_function_app.test.id
  language        = "CSharp"

  file {
    name    = "run.csx"
    content = file("testdata/run.csx")
  }

  config_json = jsonencode({
    "bindings" = [
      {
        "authLevel" = "function"
        "direction" = "in"
        "methods" = [
          "get",
          "post",
        ]
        "name" = "req"
        "type" = "httpTrigger"
      },
      {
        "direction" = "out"
        "name"      = "$return"
        "type"      = "http"
      },
    ]
  })
}
`, FunctionAppFunctionResource{}.templateWindows(data), data.RandomInteger)
}
