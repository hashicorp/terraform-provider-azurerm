package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureSshPublicKey_CreateUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ssh_key", "test")
	var d compute.SSHPublicKeyResource

	key1 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
	key2 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureSshPublicKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureSshPublicKey_template(data, key1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureSshPublicKeyExists(data.ResourceName, &d, true),
					resource.TestCheckResourceAttr(data.ResourceName, "public_key", key1),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureSshPublicKey_template(data, key2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureSshPublicKeyExists(data.ResourceName, &d, true),
					resource.TestCheckResourceAttr(data.ResourceName, "public_key", key2),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureSshPublicKeyExists(resourceName string, d *compute.SSHPublicKeyResource, shouldExist bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.SSHPublicKeysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		dName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for ssh public key: %s", dName)
		}

		resp, err := client.Get(ctx, resourceGroup, dName)
		if err != nil {
			return fmt.Errorf("Bad: Get on SSHPublicKeysClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound && shouldExist {
			return fmt.Errorf("Bad: SSH Public Key %q (resource group %q) does not exist", dName, resourceGroup)
		}
		if resp.StatusCode != http.StatusNotFound && !shouldExist {
			return fmt.Errorf("Bad: SSH Public Key %q (resource group %q) still exists", dName, resourceGroup)
		}

		*d = resp

		return nil
	}
}

func testCheckAzureSshPublicKeyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Compute.SSHPublicKeysClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_ssh_public_key" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("SSH Key still exists: \n%#v", resp.SSHPublicKeyResourceProperties)
		}
	}

	return nil
}

func testAccAzureSshPublicKey_template(data acceptance.TestData, sshKey string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
  }

resource "azurerm_ssh_public_key" "test" {
	name                = "test-public-key-%d"
	resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	public_key = "%s"
	tags = {
	  test-tag: "test-value-%d"
	}
  
  }
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sshKey, data.RandomInteger)
}
