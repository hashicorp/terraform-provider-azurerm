package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SSHPublicKeyResource struct {
}

func TestAccSshPublicKey_CreateUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ssh_public_key", "test")
	r := SSHPublicKeyResource{}

	key1 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"
	key2 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0/NDMj2wG6bSa6jbn6E3LYlUsYiWMp1CQ2sGAijPALW6OrSu30lz7nKpoh8Qdw7/A4nAJgweI5Oiiw5/BOaGENM70Go+VM8LQMSxJ4S7/8MIJEZQp5HcJZ7XDTcEwruknrd8mllEfGyFzPvJOx6QAQocFhXBW6+AlhM3gn/dvV5vdrO8ihjET2GoDUqXPYC57ZuY+/Fz6W3KV8V97BvNUhpY5yQrP5VpnyvvXNFQtzDfClTvZFPuoHQi3/KYPi6O0FSD74vo8JOBZZY09boInPejkm9fvHQqfh0bnN7B6XJoUwC1Qprrx+XIy7ust5AEn5XL7d4lOvcR14MxDDKEp you@me.com"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.template(data, key1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_key").HasValue(key1),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data, key2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_key").HasValue(key2),
			),
		},
		data.ImportStep(),
	})
}

func (t SSHPublicKeyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SSHPublicKeyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.SSHPublicKeysClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute SSH Public Key %q", id.String())
	}

	return utils.Bool(resp.ID != nil), nil
}

func (SSHPublicKeyResource) template(data acceptance.TestData, sshKey string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "ACCTESTRG-%d"
  location = "%s"
}

resource "azurerm_ssh_public_key" "test" {
  name                = "test-public-key-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  public_key          = "%s"
  tags = {
    test-tag : "test-value-%d"
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, sshKey, data.RandomInteger)
}
