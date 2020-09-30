package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementCustomDomain_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementCustomDomainExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_custom_domain", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementCustomDomain_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementCustomDomainExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementCustomDomain_requiresImport),
		},
	})
}

func testCheckAzureRMApiManagementCustomDomainDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ServiceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_custom_domain" {
			continue
		}

		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			if resp.ServiceProperties != nil && resp.ServiceProperties.HostnameConfigurations != nil && len(*resp.ServiceProperties.HostnameConfigurations) > 0 {
				return fmt.Errorf("Bad: Expected there to be no Custom Domains in the hostname_configurations field: %+v", resp.ServiceProperties.HostnameConfigurations)
			}

			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMApiManagementCustomDomainExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Custom Domains on API Management Service %q / Resource Group: %q does not exist (because API Management Service %q does not exist)", serviceName, resourceGroup, serviceName)
			}

			if resp.ServiceProperties == nil || resp.ServiceProperties.HostnameConfigurations == nil || len(*resp.ServiceProperties.HostnameConfigurations) == 0 {
				return fmt.Errorf("Bad: Expected there to be Custom Domains defined in the hostname_configurations field for API Management Service %q / Resource Group: %q", serviceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementCustomDomainsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementCustomDomain_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementCustomDomain_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "test" {
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  proxy {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, template)
}

func testAccAzureRMApiManagementCustomDomain_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementCustomDomain_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_custom_domain" "import" {
  api_management_name = azurerm_api_management_custom_domain.test.api_management_name
  resource_group_name = azurerm_api_management_custom_domain.test.resource_group_name

  proxy {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, template)
}

func testAccAzureRMApiManagementCustomDomain_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "get",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_api_management.test.identity.principal_id

    certificate_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%[3]s"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=api.example.com"
      validity_in_months = 12

      subject_alternative_names {
        dns_names = [
          "api.example.com",
          "portal.example.com",
        ]
      }
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
