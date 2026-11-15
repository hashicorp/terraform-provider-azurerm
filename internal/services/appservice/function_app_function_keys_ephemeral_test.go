// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type FunctionAppFunctionKeysEphemeral struct{}

func TestAccEphemeralFunctionAppFunctionKeys_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azurerm_function_app_function_keys", "test")
	r := FunctionAppFunctionKeysEphemeral{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0-rc1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		ProtoV6ProviderFactories: framework.ProtoV6ProviderFactoriesInit(context.Background(), "azurerm", "echo"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("default_key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("keys").AtMapKey("default"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("keys"), knownvalue.MapSizeExact(1)),
				},
			},
		},
	})
}

func (FunctionAppFunctionKeysEphemeral) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

ephemeral "azurerm_function_app_function_keys" "test" {
  function_app_id = azurerm_function_app_function.test.function_app_id
  function_name   = azurerm_function_app_function.test.name
}

provider "echo" {
  data = ephemeral.azurerm_function_app_function_keys.test
}

resource "echo" "test" {}
`, FunctionAppFunctionResource{}.withLocalFiles(data))
}
