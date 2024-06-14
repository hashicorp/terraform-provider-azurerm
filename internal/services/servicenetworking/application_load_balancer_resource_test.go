// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicenetworking_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApplicationLoadBalancerResource struct{}

func (r ApplicationLoadBalancerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := trafficcontrollerinterface.ParseTrafficControllerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceNetworking.TrafficControllerInterface.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccApplicationLoadBalancer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer", "test")

	r := ApplicationLoadBalancerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("primary_configuration_endpoint").MatchesRegex(regexp.MustCompile(`^.+\.alb.azure.com$`)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationLoadBalancer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer", "test")

	r := ApplicationLoadBalancerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("primary_configuration_endpoint").MatchesRegex(regexp.MustCompile(`^.+\.alb.azure.com$`)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationLoadBalancer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer", "test")

	r := ApplicationLoadBalancerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("primary_configuration_endpoint").MatchesRegex(regexp.MustCompile(`^.+\.alb.azure.com$`)),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("primary_configuration_endpoint").MatchesRegex(regexp.MustCompile(`^.+\.alb.azure.com$`)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationLoadBalancer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer", "test")

	r := ApplicationLoadBalancerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("primary_configuration_endpoint").MatchesRegex(regexp.MustCompile(`^.+\.alb.azure.com$`)),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r ApplicationLoadBalancerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-alb-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApplicationLoadBalancerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
  }
}

	%s

resource "azurerm_application_load_balancer" "test" {
  name                = "acctestalb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationLoadBalancerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
  }
}

	%s

resource "azurerm_application_load_balancer" "test" {
  name                = "acctestalb-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    key = "value"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationLoadBalancerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_application_load_balancer" "import" {
  name                = azurerm_application_load_balancer.test.name
  location            = azurerm_application_load_balancer.test.location
  resource_group_name = azurerm_application_load_balancer.test.resource_group_name
}
`, r.basic(data))
}
