// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

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

type KeyVaultSecretEphemeral struct{}

func TestAccEphemeralKeyVaultSecret_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azurerm_key_vault_secret", "test")
	r := KeyVaultSecretEphemeral{}

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
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("value"), knownvalue.StringExact("rick-and-morty")),
				},
			},
		},
	})
}

func TestAccEphemeralKeyVaultSecret_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azurerm_key_vault_secret", "test")
	r := KeyVaultSecretEphemeral{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0-rc1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		ProtoV6ProviderFactories: framework.ProtoV6ProviderFactoriesInit(context.Background(), "azurerm", "echo"),
		Steps: []resource.TestStep{
			{
				Config: r.complete(data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("value"), knownvalue.StringExact("<rick><morty /></rick>")),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("not_before_date"), knownvalue.StringExact("2019-01-01T01:02:03Z")),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("expiration_date"), knownvalue.StringExact("2020-01-01T01:02:03Z")),
				},
			},
		},
	})
}

func (KeyVaultSecretEphemeral) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

ephemeral "azurerm_key_vault_secret" "test" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
}

provider "echo" {
  data = ephemeral.azurerm_key_vault_secret.test
}

resource "echo" "test" {}
`, KeyVaultSecretResource{}.basic(data))
}

func (KeyVaultSecretEphemeral) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

ephemeral "azurerm_key_vault_secret" "test" {
  name         = azurerm_key_vault_secret.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_secret.test.version
}

provider "echo" {
  data = ephemeral.azurerm_key_vault_secret.test
}

resource "echo" "test" {}
`, KeyVaultSecretResource{}.complete(data))
}
