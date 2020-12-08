package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceApplicationSecurityGroupAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_multipleIPConfigurations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_multipleIPConfigurations(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceApplicationSecurityGroupAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceApplicationSecurityGroupAssociationExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_interface_application_security_group_association"),
			},
		},
	})
}

func TestAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_deleted(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceApplicationSecurityGroupAssociationExists(data.ResourceName),
					testCheckAzureRMNetworkInterfaceApplicationSecurityGroupAssociationDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_updateNIC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_interface_application_security_group_association", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		// intentional as this is a Virtual Resource
		CheckDestroy: testCheckAzureRMNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceApplicationSecurityGroupAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_updateNIC(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkInterfaceApplicationSecurityGroupAssociationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetworkInterfaceApplicationSecurityGroupAssociationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		nicID, err := azure.ParseAzureResourceID(rs.Primary.Attributes["network_interface_id"])
		if err != nil {
			return err
		}

		nicName := nicID.Path["networkInterfaces"]
		resourceGroup := nicID.ResourceGroup
		applicationSecurityGroupId := rs.Primary.Attributes["application_security_group_id"]

		read, err := client.Get(ctx, resourceGroup, nicName, "")
		if err != nil {
			return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		found := false
		for _, config := range *read.InterfacePropertiesFormat.IPConfigurations {
			if config.ApplicationSecurityGroups != nil {
				for _, group := range *config.ApplicationSecurityGroups {
					if *group.ID == applicationSecurityGroupId {
						found = true
						break
					}
				}
			}
		}

		if !found {
			return fmt.Errorf("Association between NIC %q and Application Security Group %q was not found!", nicName, applicationSecurityGroupId)
		}

		return nil
	}
}

func testCheckAzureRMNetworkInterfaceApplicationSecurityGroupAssociationDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.InterfacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		nicID, err := azure.ParseAzureResourceID(rs.Primary.Attributes["network_interface_id"])
		if err != nil {
			return err
		}

		nicName := nicID.Path["networkInterfaces"]
		resourceGroup := nicID.ResourceGroup
		applicationSecurityGroupId := rs.Primary.Attributes["application_security_group_id"]

		read, err := client.Get(ctx, resourceGroup, nicName, "")
		if err != nil {
			return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		configs := *read.InterfacePropertiesFormat.IPConfigurations
		for _, config := range configs {
			if config.ApplicationSecurityGroups != nil {
				groups := make([]network.ApplicationSecurityGroup, 0)
				for _, group := range *config.ApplicationSecurityGroups {
					if *group.ID != applicationSecurityGroupId {
						groups = append(groups, group)
					}
				}
				config.ApplicationSecurityGroups = &groups
			}
		}

		read.InterfacePropertiesFormat.IPConfigurations = &configs

		future, err := client.CreateOrUpdate(ctx, resourceGroup, nicName, read)
		if err != nil {
			return fmt.Errorf("Error removing Application Security Group Association for Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for removal of Application Security Group Association for NIC %q (Resource Group %q): %+v", nicName, resourceGroup, err)
		}

		return nil
	}
}

func testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_basic(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_application_security_group_association" "test" {
  network_interface_id          = azurerm_network_interface.test.id
  application_security_group_id = azurerm_application_security_group.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface_application_security_group_association" "import" {
  network_interface_id          = azurerm_network_interface_application_security_group_association.test.network_interface_id
  application_security_group_id = azurerm_network_interface_application_security_group_association.test.application_security_group_id
}
`, template)
}

func testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_multipleIPConfigurations(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "testconfiguration2"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_interface_application_security_group_association" "test" {
  network_interface_id          = azurerm_network_interface.test.id
  application_security_group_id = azurerm_application_security_group.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_updateNIC(data acceptance.TestData) string {
	template := testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_interface" "test" {
  name                = "acctestni-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    primary                       = true
  }

  ip_configuration {
    name                          = "testconfiguration2"
    private_ip_address_version    = "IPv6"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_network_interface_application_security_group_association" "test" {
  network_interface_id          = azurerm_network_interface.test.id
  application_security_group_id = azurerm_application_security_group.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMNetworkInterfaceApplicationSecurityGroupAssociation_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
