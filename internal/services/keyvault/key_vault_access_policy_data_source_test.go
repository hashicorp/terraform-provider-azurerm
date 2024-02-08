// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KeyVaultAccessPolicyDataSource struct{}

func TestAccDataSourceKeyVaultAccessPolicy_key(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.testAccDataSourceKeyVaultAccessPolicy("Key Management"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_permissions.#").HasValue("12"),
				check.That(data.ResourceName).Key("secret_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate_permissions.#").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultAccessPolicy_secret(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.testAccDataSourceKeyVaultAccessPolicy("Secret Management"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("secret_permissions.#").HasValue("7"),
				check.That(data.ResourceName).Key("certificate_permissions.#").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultAccessPolicy_certificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.testAccDataSourceKeyVaultAccessPolicy("Certificate Management"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("secret_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate_permissions.#").HasValue("15"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultAccessPolicy_keySecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.testAccDataSourceKeyVaultAccessPolicy("Key & Secret Management"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_permissions.#").HasValue("12"),
				check.That(data.ResourceName).Key("secret_permissions.#").HasValue("7"),
				check.That(data.ResourceName).Key("certificate_permissions.#").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultAccessPolicy_keyCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.testAccDataSourceKeyVaultAccessPolicy("Key & Certificate Management"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_permissions.#").HasValue("12"),
				check.That(data.ResourceName).Key("secret_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("certificate_permissions.#").HasValue("15"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultAccessPolicy_secretCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.testAccDataSourceKeyVaultAccessPolicy("Secret & Certificate Management"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_permissions.#").HasValue("0"),
				check.That(data.ResourceName).Key("secret_permissions.#").HasValue("7"),
				check.That(data.ResourceName).Key("certificate_permissions.#").HasValue("15"),
			),
		},
	})
}

func TestAccDataSourceKeyVaultAccessPolicy_keySecretCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	r := KeyVaultAccessPolicyDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.testAccDataSourceKeyVaultAccessPolicy("Key, Secret, & Certificate Management"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_permissions.#").HasValue("12"),
				check.That(data.ResourceName).Key("secret_permissions.#").HasValue("7"),
				check.That(data.ResourceName).Key("certificate_permissions.#").HasValue("15"),
			),
		},
	})
}

func (r KeyVaultAccessPolicyDataSource) testAccDataSourceKeyVaultAccessPolicy(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_key_vault_access_policy" "test" {
  name = "%s"
}
`, name)
}
