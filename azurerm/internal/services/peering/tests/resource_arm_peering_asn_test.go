package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/peering/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPeerAsn_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_peer_asn", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPeerAsnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPeerAsn_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPeerAsnExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPeerAsn_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_peer_asn", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPeerAsnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPeerAsn_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPeerAsnExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPeerAsn_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_peer_asn", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPeerAsnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPeerAsn_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPeerAsnExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPeerAsn_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPeerAsnExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPeerAsn_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPeerAsnExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPeerAsn_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_peer_asn", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPeerAsnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPeerAsn_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPeerAsnExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPeerAsn_requiresImport),
		},
	})
}

func testCheckAzureRMPeerAsnExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Peering.PeerAsnsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Peer Asn not found: %s", resourceName)
		}

		id, err := parse.PeerAsnID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Peer Asn %q does not exist", id.Name)
			}
			return fmt.Errorf("failed to get on Peering.PeerAsns: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPeerAsnDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Peering.PeerAsnsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_peer_asn" {
			continue
		}

		id, err := parse.PeerAsnID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("failed to get on Peering.PeerAsns: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPeerAsn_basic(data acceptance.TestData) string {
	template := testAccAzureRMPeerAsn_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_peer_asn" "test" {
  name = "acctestASN_%d"
  asn  = 123
  contact {
    role  = "Noc"
    email = "email@test.com"
  }
  peer_name = "acctest-peer"
}
`, template, data.RandomInteger)
}

func testAccAzureRMPeerAsn_complete(data acceptance.TestData) string {
	template := testAccAzureRMPeerAsn_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_peer_asn" "test" {
  name = "acctestASN_%d"
  asn  = 123
  contact {
    role  = "Noc"
    email = "email@test.com"
    phone = 12345
  }
  contact {
    role  = "Service"
    email = "email@test.com"
    phone = 12345
  }
  peer_name = "acctest-peer2"
}
`, template, data.RandomInteger)
}

func testAccAzureRMPeerAsn_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPeerAsn_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_peer_asn" "import" {
  name = azurerm_peer_asn.test.name
  asn  = azurerm_peer_asn.test.asn
  contact {
    role  = "Noc"
    email = "email@test.com"
  }
  peer_name = azurerm_peer_asn.test.peer_name
}
`, template)
}

func testAccAzureRMPeerAsn_template(data acceptance.TestData) string {
	return `provider "azurerm" {
  features {}
}
`
}
