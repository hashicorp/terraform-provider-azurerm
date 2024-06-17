// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicenetworking_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/associationsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SubnetAssociationResource struct{}

func (r SubnetAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := associationsinterface.ParseAssociationID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	resp, err := clients.ServiceNetworking.AssociationsInterface.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("while checking existence of %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccApplicationLoadBalancerSubnetAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_subnet_association", "test")

	r := SubnetAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationLoadBalancerSubnetAssociation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_subnet_association", "test")

	r := SubnetAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationLoadBalancerSubnetAssociation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_subnet_association", "test")

	r := SubnetAssociationResource{}
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

func TestAccApplicationLoadBalancerSubnetAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_subnet_association", "test")

	r := SubnetAssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r SubnetAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-alb-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_load_balancer" "test" {
  name                = "acctestalb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.ServiceNetworking/trafficControllers"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r SubnetAssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
  }
}

%s

resource "azurerm_application_load_balancer_subnet_association" "test" {
  name                         = "acct-%d"
  application_load_balancer_id = azurerm_application_load_balancer.test.id
  subnet_id                    = azurerm_subnet.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r SubnetAssociationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
  }
}

%s

resource "azurerm_application_load_balancer_subnet_association" "test" {
  name                         = "acct-%d"
  application_load_balancer_id = azurerm_application_load_balancer.test.id
  subnet_id                    = azurerm_subnet.test.id
  tags = {
    key = "value"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SubnetAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_application_load_balancer_subnet_association" "import" {
  name                         = azurerm_application_load_balancer_subnet_association.test.name
  application_load_balancer_id = azurerm_application_load_balancer_subnet_association.test.application_load_balancer_id
  subnet_id                    = azurerm_application_load_balancer_subnet_association.test.subnet_id
}
`, r.basic(data))
}
