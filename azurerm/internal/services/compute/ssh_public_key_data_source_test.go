package compute

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SSHPublicKeyDataSource struct {
}

func TestAccDataSourceAzureSshPublicKey_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_ssh_public_key", "test")
	r := SSHPublicKeyDataSource{}

	key1 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.template(data, key1),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("public_key").HasValue(key1),
			),
		},
	})
}

func (SSHPublicKeyDataSource) template(data acceptance.TestData, sshKey string) string {
	return fmt.Sprintf(`
%s

data "azurerm_ssh_public_key" "test" {
  name                = azurerm_ssh_public_key.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, SSHPublicKeyResource{}.template(data, sshKey))
}
