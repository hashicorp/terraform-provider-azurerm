package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_key(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_access_policy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Key Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "key_permissions.#", "9"),
					resource.TestCheckNoResourceAttr(dataSourceName, "secret_permissions"),
					resource.TestCheckNoResourceAttr(dataSourceName, "certificate_permissions"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_secret(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_access_policy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Secret Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(dataSourceName, "key_permissions"),
					resource.TestCheckResourceAttr(dataSourceName, "secret_permissions.#", "7"),
					resource.TestCheckNoResourceAttr(dataSourceName, "certificate_permissions"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_certificate(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_access_policy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Certificate Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(dataSourceName, "key_permissions"),
					resource.TestCheckNoResourceAttr(dataSourceName, "secret_permissions"),
					resource.TestCheckResourceAttr(dataSourceName, "certificate_permissions.#", "12"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_keySecret(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_access_policy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Key & Secret Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "key_permissions.#", "9"),
					resource.TestCheckResourceAttr(dataSourceName, "secret_permissions.#", "7"),
					resource.TestCheckNoResourceAttr(dataSourceName, "certificate_permissions"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_keyCertificate(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_access_policy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Key & Certificate Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "key_permissions.#", "9"),
					resource.TestCheckNoResourceAttr(dataSourceName, "secret_permissions"),
					resource.TestCheckResourceAttr(dataSourceName, "certificate_permissions.#", "12"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_secretCertificate(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_access_policy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Secret & Certificate Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(dataSourceName, "key_permissions"),
					resource.TestCheckResourceAttr(dataSourceName, "secret_permissions.#", "7"),
					resource.TestCheckResourceAttr(dataSourceName, "certificate_permissions.#", "12"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKeyVaultAccessPolicy_keySecretCertificate(t *testing.T) {
	dataSourceName := "data.azurerm_key_vault_access_policy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeyVaultAccessPolicy("Key, Secret, & Certificate Management"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "key_permissions.#", "9"),
					resource.TestCheckResourceAttr(dataSourceName, "secret_permissions.#", "7"),
					resource.TestCheckResourceAttr(dataSourceName, "certificate_permissions.#", "12"),
				),
			},
		},
	})
}

func testAccDataSourceKeyVaultAccessPolicy(name string) string {
	return fmt.Sprintf(`
data "azurerm_key_vault_access_policy" "test" {
  name = "%s"
}
`, name)
}
