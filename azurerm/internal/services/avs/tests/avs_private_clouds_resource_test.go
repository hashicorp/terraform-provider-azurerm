package tests

import (
    "fmt"
    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/terraform"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/avs/parse"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
    "testing"
)

func TestAccAzureRMavsPrivateCloud_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { acceptance.PreCheck(t) },
        Providers:    acceptance.SupportedProviders,
        CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccAzureRMavsPrivateCloud_basic(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.ImportStep(),
        },
    })
}

func TestAccAzureRMavsPrivateCloud_requiresImport(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { acceptance.PreCheck(t) },
        Providers:    acceptance.SupportedProviders,
        CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccAzureRMavsPrivateCloud_basic(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.RequiresImportErrorStep(testAccAzureRMavsPrivateCloud_requiresImport),
        },
    })
}

func TestAccAzureRMavsPrivateCloud_complete(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { acceptance.PreCheck(t) },
        Providers:    acceptance.SupportedProviders,
        CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccAzureRMavsPrivateCloud_complete(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.ImportStep(),
        },
    })
}

func TestAccAzureRMavsPrivateCloud_update(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { acceptance.PreCheck(t) },
        Providers:    acceptance.SupportedProviders,
        CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccAzureRMavsPrivateCloud_basic(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.ImportStep(),
            {
                Config: testAccAzureRMavsPrivateCloud_complete(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.ImportStep(),
            {
                Config: testAccAzureRMavsPrivateCloud_basic(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.ImportStep(),
        },
    })
}

func TestAccAzureRMavsPrivateCloud_updateManagementCluster(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { acceptance.PreCheck(t) },
        Providers:    acceptance.SupportedProviders,
        CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccAzureRMavsPrivateCloud_complete(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.ImportStep(),
            {
                Config: testAccAzureRMavsPrivateCloud_updateManagementCluster(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.ImportStep(),
        },
    })
}

func TestAccAzureRMavsPrivateCloud_updateIdentitySources(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_avs_private_cloud", "test")
    resource.ParallelTest(t, resource.TestCase{
        PreCheck:     func() { acceptance.PreCheck(t) },
        Providers:    acceptance.SupportedProviders,
        CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccAzureRMavsPrivateCloud_complete(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.ImportStep(),
            {
                Config: testAccAzureRMavsPrivateCloud_updateIdentitySources(data),
                Check: resource.ComposeTestCheckFunc(
                    testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
                ),
            },
            data.ImportStep(),
        },
    })
}

func testCheckAzureRMavsPrivateCloudExists(resourceName string) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        client := acceptance.AzureProvider.Meta().(*clients.Client).Avs.PrivateCloudClient
        ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
        rs, ok := s.RootModule().Resources[resourceName]
        if !ok {
            return fmt.Errorf("avs PrivateCloud not found: %s", resourceName)
        }
        id, err := parse.AvsPrivateCloudID(rs.Primary.ID)
        if err != nil {
            return err
        }
        if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
            if !utils.ResponseWasNotFound(resp.Response) {
                return fmt.Errorf("bad: Avs PrivateCloud %q does not exist", id.Name)
            }
            return fmt.Errorf("bad: Get on Avs.PrivateCloudClient: %+v", err)
        }
        return nil
    }
}

func testCheckAzureRMavsPrivateCloudDestroy(s *terraform.State) error {
    client := acceptance.AzureProvider.Meta().(*clients.Client).Avs.PrivateCloudClient
    ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "azurerm_avs_private_cloud" {
            continue
        }
        id, err := parse.AvsPrivateCloudID(rs.Primary.ID)
        if err != nil {
            return err
        }
        if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
            if !utils.ResponseWasNotFound(resp.Response) {
                 return fmt.Errorf("bad: Get on Avs.PrivateCloudClient: %+v", err)
            }
        }
        return nil
    }
    return nil
}

func testAccAzureRMavsPrivateCloud_template(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-avs-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMavsPrivateCloud_basic(data acceptance.TestData) string {
    template := testAccAzureRMavsPrivateCloud_template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name = "acctest-apc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  sku {
    name = "av36"
  }

  management_cluster {
    cluster_size = 4
  }
  network_block = "192.168.48.0/22"
}
`, template, data.RandomInteger)
}

func testAccAzureRMavsPrivateCloud_requiresImport(data acceptance.TestData) string {
    config := testAccAzureRMavsPrivateCloud_basic(data)
    return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "import" {
  name = azurerm_avs_private_cloud.test.name
  resource_group_name = azurerm_avs_private_cloud.test.resource_group_name
  location = azurerm_avs_private_cloud.test.location
  sku {
    name = azurerm_avs_private_cloud.test.sku.name
  }
  management_cluster {
    cluster_size = azurerm_avs_private_cloud.test.management_cluster.cluster_size
  }
  network_block = azurerm_avs_private_cloud.test.network_block
}
`, config)
}

func testAccAzureRMavsPrivateCloud_complete(data acceptance.TestData) string {
    template := testAccAzureRMavsPrivateCloud_template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name = "acctest-apc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  sku {
    name = "AV36"
  }

  management_cluster {
    cluster_size = 4
  }
  network_block = "192.168.48.0/22"
  identity_source {
    name = ""
    alias = ""
    base_group_dn = ""
    base_user_dn = ""
    domain = ""
    password = ""
    primary_server = ""
    secondary_server = ""
    ssl = false
    username = ""
  }
  internet = false
  nsxt_password = ""
  vcenter_password = ""
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMavsPrivateCloud_updateManagementCluster(data acceptance.TestData) string {
    template := testAccAzureRMavsPrivateCloud_template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name = "acctest-apc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  sku {
    name = "AV36"
  }

  management_cluster {
    cluster_size = 4
  }
  network_block = "192.168.48.0/22"
  identity_source {
    name = ""
    alias = ""
    base_group_dn = ""
    base_user_dn = ""
    domain = ""
    password = ""
    primary_server = ""
    secondary_server = ""
    ssl = false
    username = ""
  }
  internet = false
  nsxt_password = ""
  vcenter_password = ""
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMavsPrivateCloud_updateIdentitySources(data acceptance.TestData) string {
    template := testAccAzureRMavsPrivateCloud_template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_avs_private_cloud" "test" {
  name = "acctest-apc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  sku {
    name = "AV36"
  }

  management_cluster {
    cluster_size = 4
  }
  network_block = "192.168.48.0/22"
  identity_source {
    name = ""
    alias = ""
    base_group_dn = ""
    base_user_dn = ""
    domain = ""
    password = ""
    primary_server = ""
    secondary_server = ""
    ssl = false
    username = ""
  }
  internet = false
  nsxt_password = ""
  vcenter_password = ""
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
