package redhatopenshift_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type OpenShiftClusterResource struct {
}

var (
	clientId     = os.Getenv("ARM_CLIENT_ID")
	clientSecret = os.Getenv("ARM_CLIENT_SECRET")
)

func (t OpenShiftClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.RedHatOpenshift.OpenShiftClustersClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		return nil, fmt.Errorf("reading Red Hat Openshift Cluster (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func TestAccOpenShiftCluster_version(t *testing.T) {
	testAccOpenShiftCluster_version(t)
}

func testAccOpenShiftCluster_version(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhatopenshift_cluster", "test")
	r := OpenShiftClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.version(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("version").Exists(),
			),
		},
	})
}

func (OpenShiftClusterResource) version(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aro-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/22"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "master_subnet" {
  name                 = "master-subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/23"]
  service_endpoints    = ["Microsoft.ContainerRegistry"]
  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_subnet" "worker_subnet" {
  name                 = "worker-subnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/23"]
  service_endpoints    = ["Microsoft.ContainerRegistry"]
}

resource "azurerm_redhatopenshift_cluster" "test" {
  name                = "acctestaro%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  
  master_profile {
		vm_size 	= "Standard_D8s_v3"
    subnet_id = azurerm_subnet.master_subnet.id
  }
  
  worker_profile {
		vm_size   = "Standard_D4s_v3"
    subnet_id = azurerm_subnet.worker_subnet.id
  }

	service_principal {
    client_id     = %q
    client_secret = %q
  }
}
  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, clientId, clientSecret)
}
