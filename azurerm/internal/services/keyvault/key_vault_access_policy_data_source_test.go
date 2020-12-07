package keyvault

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_key(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Key Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.#", "9"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "secret_permissions"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "certificate_permissions"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_secret(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Secret Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(data.ResourceName, "key_permissions"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.#", "7"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "certificate_permissions"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_certificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Certificate Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(data.ResourceName, "key_permissions"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "secret_permissions"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_permissions.#", "12"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_keySecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Key & Secret Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.#", "9"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.#", "7"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "certificate_permissions"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_keyCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Key & Certificate Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.#", "9"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "secret_permissions"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_permissions.#", "12"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_secretCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Secret & Certificate Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(data.ResourceName, "key_permissions"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.#", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_permissions.#", "12"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_keySecretCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_key_vault_access_policy", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Key, Secret, & Certificate Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "key_permissions.#", "9"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret_permissions.#", "7"),
					resource.TestCheckResourceAttr(data.ResourceName, "certificate_permissions.#", "12"),
				),
			},
		},
	})
}

func testAccDataSourceKeyVaultAccessPolicy(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_key_vault_access_policy" "test" {
  name = "%s"
}
`, name)
}
