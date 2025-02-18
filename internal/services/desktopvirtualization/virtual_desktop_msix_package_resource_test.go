package desktopvirtualization_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/msixpackage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualDesktopMSIXPackageTestResource struct{}

func TestAccVirtualDesktopMSIXPackage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_msix_package", "test")
	r := VirtualDesktopMSIXPackageTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualDesktopMSIXPackage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_msix_package", "test")
	r := VirtualDesktopMSIXPackageTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccVirtualDesktopMSIXPackage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_msix_package", "test")
	r := VirtualDesktopMSIXPackageTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualDesktopMSIXPackage_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_desktop_msix_package", "test")
	r := VirtualDesktopMSIXPackageTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (VirtualDesktopMSIXPackageTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := msixpackage.ParseMsixPackageID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.DesktopVirtualization.MSIXPackageClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (VirtualDesktopMSIXPackageTestResource) basic(data acceptance.TestData, withPackageDependencies bool) string {
	packageDependencies := ""
	if withPackageDependencies {
		packageDependencies = `
  package_dependency {
    dependency_name = "dependency-1"
    min_version     = "0.0.0.2"
    publisher       = "publisher-1"
  }
`
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                = "acctestHP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  type                = "Pooled"
  load_balancer_type  = "BreadthFirst"
}

resource "azurerm_virtual_desktop_msix_package" "test" {
  name                  = "acctestMSIXPackage"
  host_pool_name        = azurerm_virtual_desktop_host_pool.test.name
  resource_group_name   = azurerm_resource_group.test.name
  image_path            = "//path/to/image.vhd"
  last_updated_in_utc   = "2021-09-01T00:00:00"

  package_application {
    app_id            = "app-1"
    app_user_model_id = "app-user-model-1"
    description       = "Testing app 1"
    friendly_name     = "I am your friendly neighbourhood testing app 1"
    icon_image_name   = "icon.png"
    raw_icon          = "VGhpcyBpcyBhIHN0cmluZyB0byBoYXNo"
    raw_png           = "VGhpcyBpcyBhIHN0cmluZyB0byBwYWdl"
  }

  package_family_name   = "msix-package-family-1"
  package_name          = "msix-package-1"
  package_relative_path = "path/to/package"
  version               = "0.0.0.1"

  %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, packageDependencies)
}

func (r VirtualDesktopMSIXPackageTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_desktop_msix_package" "import" {
  name                  = "acctestMSIXPackage"
  host_pool_name        = azurerm_virtual_desktop_host_pool.test.name
  resource_group_name   = azurerm_resource_group.test.name
  image_path            = "//path/to/image.vhd"
  last_updated_in_utc   = "2021-09-01T00:00:00"

  package_application {
    app_id            = "app-1"
    app_user_model_id = "app-user-model-1"
    description       = "Testing app 1"
    friendly_name     = "I am your friendly neighbourhood testing app 1"
    icon_image_name   = "icon.png"
    raw_icon          = "VGhpcyBpcyBhIHN0cmluZyB0byBoYXNo"
    raw_png           = "VGhpcyBpcyBhIHN0cmluZyB0byBwYWdl"
  }

  package_family_name   = "msix-package-family-1"
  package_name          = "msix-package-1"
  package_relative_path = "path/to/package"
  version               = "0.0.0.1"
}
`, r.basic(data, false))
}

func (r VirtualDesktopMSIXPackageTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vdesktop-%d"
  location = "%s"
}

resource "azurerm_virtual_desktop_host_pool" "test" {
  name                = "acctestHP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  type                = "Pooled"
  load_balancer_type  = "BreadthFirst"
}

resource "azurerm_virtual_desktop_msix_package" "test" {
  name                  = "acctestMSIXPackage"
  host_pool_name        = azurerm_virtual_desktop_host_pool.test.name
  resource_group_name   = azurerm_resource_group.test.name
  image_path            = "//path/to/image.vhd"
  last_updated_in_utc   = "2021-09-01T00:00:00"

  package_application {
    app_id            = "app-1"
    app_user_model_id = "app-user-model-1"
    description       = "Testing app 1"
    friendly_name     = "I am your friendly neighbourhood testing app 1"
    icon_image_name   = "icon.png"
    raw_icon          = "VGhpcyBpcyBhIHN0cmluZyB0byBoYXNo"
    raw_png           = "VGhpcyBpcyBhIHN0cmluZyB0byBwYWdl"
  }

  package_family_name   = "msix-package-family-1"
  package_name          = "msix-package-1"
  package_relative_path = "path/to/package"
  version               = "0.0.0.1"

  display_name = "acctestMSIXPackage"
  enabled = true
  
  package_dependency {
    dependency_name = "dependency-1"
    min_version     = "0.0.0.2"
    publisher       = "publisher-1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
