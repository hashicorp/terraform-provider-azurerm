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

type KeyVaultCertificateEphemeral struct{}

func TestAccEphemeralKeyVaultCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateEphemeral{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0-rc1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		ProtoV6ProviderFactories: framework.ProtoV6ProviderFactoriesInit(context.Background(), "azurerm", "echo"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("hex"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("pem"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("not_before_date"), knownvalue.StringExact("2017-10-10T08:27:55Z")),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("expiration_date"), knownvalue.StringExact("2027-10-08T08:27:55Z")),
				},
			},
		},
	})
}

func TestAccEphemeralKeyVaultCertificate_ecdsaPFX(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateEphemeral{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0-rc1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		ProtoV6ProviderFactories: framework.ProtoV6ProviderFactoriesInit(context.Background(), "azurerm", "echo"),
		Steps: []resource.TestStep{
			{
				Config: r.format(KeyVaultCertificateResource{}.basicImportPFX_ECDSA(data)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("hex"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("pem"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("key"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccEphemeralKeyVaultCertificate_ecdsaPEM(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateEphemeral{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0-rc1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		ProtoV6ProviderFactories: framework.ProtoV6ProviderFactoriesInit(context.Background(), "azurerm", "echo"),
		Steps: []resource.TestStep{
			{
				Config: r.format(KeyVaultCertificateResource{}.basicImportPEM_ECDSA(data)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("hex"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("pem"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("key"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccEphemeralKeyVaultCertificate_rsaBundlePEM(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateEphemeral{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0-rc1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		ProtoV6ProviderFactories: framework.ProtoV6ProviderFactoriesInit(context.Background(), "azurerm", "echo"),
		Steps: []resource.TestStep{
			{
				Config: r.format(KeyVaultCertificateResource{}.basicImportPEM_RSA_bundle(data)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("hex"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("pem"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("certificate_count"), knownvalue.Int64Exact(2)),
				},
			},
		},
	})
}

func TestAccEphemeralKeyVaultCertificate_rsaSinglePEM(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateEphemeral{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0-rc1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		ProtoV6ProviderFactories: framework.ProtoV6ProviderFactoriesInit(context.Background(), "azurerm", "echo"),
		Steps: []resource.TestStep{
			{
				Config: r.format(KeyVaultCertificateResource{}.basicImportPEM_RSA(data)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("hex"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("pem"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("certificate_count"), knownvalue.Int64Exact(1)),
				},
			},
		},
	})
}

func TestAccEphemeralKeyVaultCertificate_rsaBundlePFX(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azurerm_key_vault_certificate", "test")
	r := KeyVaultCertificateEphemeral{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0-rc1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		ProtoV6ProviderFactories: framework.ProtoV6ProviderFactoriesInit(context.Background(), "azurerm", "echo"),
		Steps: []resource.TestStep{
			{
				Config: r.format(KeyVaultCertificateResource{}.basicImportPFX_RSA_bundle(data)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("hex"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("pem"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("key"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("certificate_count"), knownvalue.Int64Exact(2)),
				},
			},
		},
	})
}

func (KeyVaultCertificateEphemeral) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

ephemeral "azurerm_key_vault_certificate" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
}

provider "echo" {
  data = ephemeral.azurerm_key_vault_certificate.test
}

resource "echo" "test" {}
`, KeyVaultCertificateResource{}.basicImportPFX(data))
}

func (KeyVaultCertificateEphemeral) format(formatTemplate string) string {
	return fmt.Sprintf(`
%s

ephemeral "azurerm_key_vault_certificate" "test" {
  name         = azurerm_key_vault_certificate.test.name
  key_vault_id = azurerm_key_vault.test.id
  version      = azurerm_key_vault_certificate.test.version
}

provider "echo" {
  data = ephemeral.azurerm_key_vault_certificate.test
}

resource "echo" "test" {}
`, formatTemplate)
}
