package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagement_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagement_customProps(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_customProps(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "protocols.0.enable_http2", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate", // not returned from API, sensitive
					"hostname_configuration.0.portal.0.certificate",                    // not returned from API, sensitive
					"hostname_configuration.0.portal.0.certificate_password",           // not returned from API, sensitive
					"hostname_configuration.0.developer_portal.0.certificate",          // not returned from API, sensitive
					"hostname_configuration.0.developer_portal.0.certificate_password", // not returned from API, sensitive
					"hostname_configuration.0.proxy.0.certificate",                     // not returned from API, sensitive
					"hostname_configuration.0.proxy.0.certificate_password",            // not returned from API, sensitive
					"hostname_configuration.0.proxy.1.certificate",                     // not returned from API, sensitive
					"hostname_configuration.0.proxy.1.certificate_password",            // not returned from API, sensitive
				},
			},
		},
	})
}

func TestAccAzureRMApiManagement_signInSignUpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_signInSignUpSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_policy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_policyXmlContent(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_policyXmlLink(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"policy.0.xml_link",
				},
			},
			{
				Config: testAccAzureRMApiManagement_policyRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_virtualNetworkInternal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_virtualNetworkInternal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_type", "Internal"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_addresses.#"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_virtualNetworkInternalUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_virtualNetworkInternal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_virtualNetworkInternalAdditionalLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_virtualNetworkInternalAdditionalLocation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_type", "Internal"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "additional_location.0.private_ip_addresses.#"),
				),
			},
			data.ImportStep(),
		},
	})
}

// Api Management doesn't support hostname keyvault using UserAssigned Identity
// There will be a inevitable dependency cycle here when using SystemAssigned Identity
// 1. create SystemAssigned Identity, grant the identity certificate access
// 2. Update the hostname configuration of the keyvault certificate
func TestAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsVersionedKeyVaultId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsKeyVaultId(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsVersionedKeyVaultIdUpdateCD(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsVersionlessKeyVaultId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsKeyVaultId(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsVersionlessKeyVaultIdUpdateCD(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_consumption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_consumption(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMApiManagementDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ServiceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMApiManagementExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		apiMangementName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Api Management: %s", apiMangementName)
		}

		resp, err := conn.Get(ctx, resourceGroup, apiMangementName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Api Management %q (resource group: %q) does not exist", apiMangementName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementClient: %+v", err)
		}

		return nil
	}
}

func TestAccAzureRMApiManagement_identityUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identityUserAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identityNoneUpdateUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identityNone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identityUserAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identityUserAssignedUpdateNone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identityUserAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identityNone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identitySystemAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identitySystemAssignedUpdateNone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identitySystemAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identityNone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identityNoneUpdateSystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identityNone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identitySystemAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identitySystemAssignedUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identitySystemAssignedUserAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identitySystemAssignedUserAssignedUpdateNone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identitySystemAssignedUserAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identityNone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identityNoneUpdateSystemAssignedUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identityNone(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identitySystemAssignedUserAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identitySystemAssignedUserAssignedUpdateSystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identitySystemAssignedUserAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identitySystemAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_identitySystemAssignedUserAssignedUpdateUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_identitySystemAssignedUserAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_identityUserAssigned(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMApiManagement_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_policyXmlContent(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  policy {
    xml_content = <<XML
<policies>
  <inbound>
    <set-variable name="abc" value="@(context.Request.Headers.GetValueOrDefault("X-Header-Name", ""))" />
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>
XML

  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_policyXmlLink(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  policy {
    xml_link = "https://gist.githubusercontent.com/tombuildsstuff/4f58581599d2c9f64b236f505a361a67/raw/0d29dcb0167af1e5afe4bd52a6d7f69ba1e05e1f/example.xml"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_policyRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  policy = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagement_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "import" {
  name                = azurerm_api_management.test.name
  location            = azurerm_api_management.test.location
  resource_group_name = azurerm_api_management.test.resource_group_name
  publisher_name      = azurerm_api_management.test.publisher_name
  publisher_email     = azurerm_api_management.test.publisher_email

  sku_name = "Developer_1"
}
`, template)
}

func testAccAzureRMApiManagement_customProps(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  security {
    enable_frontend_tls10     = true
    enable_triple_des_ciphers = true
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func testAccAzureRMApiManagement_signInSignUpSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  sign_in {
    enabled = true
  }

  sign_up {
    enabled = true

    terms_of_service {
      enabled          = true
      consent_required = false
      text             = "Lorem Ipsum Dolor Morty"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-api1-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-api2-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test3" {
  name     = "acctestRG-api3-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                      = "acctestAM-%d"
  publisher_name            = "pub1"
  publisher_email           = "pub1@email.com"
  notification_sender_email = "notification@email.com"

  additional_location {
    location = azurerm_resource_group.test2.location
  }

  additional_location {
    location = azurerm_resource_group.test3.location
  }

  certificate {
    encoded_certificate  = filebase64("testdata/api_management_api_test.pfx")
    certificate_password = "terraform"
    store_name           = "CertificateAuthority"
  }

  certificate {
    encoded_certificate  = filebase64("testdata/api_management_api_test.pfx")
    certificate_password = "terraform"
    store_name           = "Root"
  }

  protocols {
    enable_http2 = true
  }

  security {
    enable_backend_tls11      = true
    enable_backend_ssl30      = true
    enable_backend_tls10      = true
    enable_frontend_ssl30     = true
    enable_frontend_tls10     = true
    enable_frontend_tls11     = true
    enable_triple_des_ciphers = true
  }

  hostname_configuration {
    proxy {
      host_name                    = "api.terraform.io"
      certificate                  = filebase64("testdata/api_management_api_test.pfx")
      certificate_password         = "terraform"
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }

    proxy {
      host_name                    = "api2.terraform.io"
      certificate                  = filebase64("testdata/api_management_api2_test.pfx")
      certificate_password         = "terraform"
      negotiate_client_certificate = true
    }

    portal {
      host_name            = "portal.terraform.io"
      certificate          = filebase64("testdata/api_management_portal_test.pfx")
      certificate_password = "terraform"
    }

    developer_portal {
      host_name   = "developer-portal.terraform.io"
      certificate = filebase64("testdata/api_management_developer_portal_test.pfx")
    }
  }

  sku_name = "Premium_1"

  tags = {
    "Acceptance" = "Test"
  }

  location            = azurerm_resource_group.test1.location
  resource_group_name = azurerm_resource_group.test1.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.Locations.Ternary, data.RandomInteger)
}

func testAccAzureRMApiManagement_virtualNetworkTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNET-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSNET-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_network_security_rule" "port_3443" {
  name                        = "Port_3443"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "3443"
  source_address_prefix       = "ApiManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test.name
  network_security_group_name = azurerm_network_security_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMApiManagement_virtualNetworkInternal(data acceptance.TestData) string {
	template := testAccAzureRMApiManagement_virtualNetworkTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  virtual_network_type = "Internal"
  virtual_network_configuration {
    subnet_id = azurerm_subnet.test.id
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagement_virtualNetworkInternalAdditionalLocation(data acceptance.TestData) string {
	template := testAccAzureRMApiManagement_virtualNetworkTemplate(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-%[2]d"
  location = "%[3]s"
}

// subnet2 from the second location
resource "azurerm_virtual_network" "test2" {
  name                = "acctestVNET2-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test2" {
  name                 = "acctestSNET2-%[2]d"
  resource_group_name  = azurerm_resource_group.test2.name
  virtual_network_name = azurerm_virtual_network.test2.name
  address_prefix       = "10.1.1.0/24"
}

resource "azurerm_network_security_group" "test2" {
  name                = "acctest-NSG2-%[2]d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
}

resource "azurerm_subnet_network_security_group_association" "test2" {
  subnet_id                 = azurerm_subnet.test2.id
  network_security_group_id = azurerm_network_security_group.test2.id
}

resource "azurerm_network_security_rule" "port_3443_2" {
  name                        = "Port_3443"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "3443"
  source_address_prefix       = "ApiManagement"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = azurerm_resource_group.test2.name
  network_security_group_name = azurerm_network_security_group.test2.name
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Premium_1"

  additional_location {
    location = azurerm_resource_group.test2.location
    virtual_network_configuration {
      subnet_id = azurerm_subnet.test2.id
    }
  }

  virtual_network_type = "Internal"
  virtual_network_configuration {
    subnet_id = azurerm_subnet.test.id
  }
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func testAccAzureRMApiManagement_identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management" "test" {
  depends_on          = [azurerm_user_assigned_identity.test]
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagement_identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagement_identityNone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
data "azurerm_client_config" "current" {}
resource "azurerm_key_vault" "test" {
  name                = "acctestKV-%[4]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}
resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id
  certificate_permissions = [
    "Create",
    "Delete",
    "Deleteissuers",
    "Get",
    "Getissuers",
    "Import",
    "List",
    "Listissuers",
    "Managecontacts",
    "Manageissuers",
    "Setissuers",
    "Update",
  ]
  secret_permissions = [
    "Delete",
    "Get",
    "List",
    "Purge",
  ]
}
resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = azurerm_api_management.test.identity[0].tenant_id
  object_id    = azurerm_api_management.test.identity[0].principal_id
  secret_permissions = [
    "Get",
    "List",
  ]
}
resource "azurerm_key_vault_certificate" "test" {
  depends_on   = [azurerm_key_vault_access_policy.test]
  name         = "acctestKVCert-%[3]d"
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
    secret_properties {
      content_type = "application/x-pkcs12"
    }
    x509_certificate_properties {
      # Server Authentication = 1.3.6.1.5.5.7.3.1
      # Client Authentication = 1.3.6.1.5.5.7.3.2
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject_alternative_names {
        dns_names = ["api.terraform.io"]
      }
      subject            = "CN=api.terraform.io"
      validity_in_months = 1
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
func testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsKeyVaultId(data acceptance.TestData) string {
	template := testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsVersionlessKeyVaultIdUpdateCD(data acceptance.TestData) string {
	template := testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
  identity {
    type = "SystemAssigned"
  }
  hostname_configuration {
    proxy {
      host_name                    = "api.terraform.io"
      key_vault_id                 = "${azurerm_key_vault.test.vault_uri}secrets/${azurerm_key_vault_certificate.test.name}"
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsVersionedKeyVaultIdUpdateCD(data acceptance.TestData) string {
	template := testAccAzureRMApiManagement_identitySystemAssignedUpdateHostnameConfigurationsTemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
  identity {
    type = "SystemAssigned"
  }
  hostname_configuration {
    proxy {
      host_name                    = "api.terraform.io"
      key_vault_id                 = azurerm_key_vault_certificate.test.secret_id
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagement_consumption(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

