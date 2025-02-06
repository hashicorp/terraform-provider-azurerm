// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicenetworking_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/frontendsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApplicationLoadBalancerFrontendResource struct{}

func (r ApplicationLoadBalancerFrontendResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := frontendsinterface.ParseFrontendID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	resp, err := clients.ServiceNetworking.FrontendsInterface.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("while checking existence for %q: %+v", id.String(), err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccApplicationLoadBalancerFrontend_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_frontend", "test")

	r := ApplicationLoadBalancerFrontendResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fully_qualified_domain_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationLoadBalancerFrontend_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_frontend", "test")

	r := ApplicationLoadBalancerFrontendResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fully_qualified_domain_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationLoadBalancerFrontend_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_frontend", "test")

	r := ApplicationLoadBalancerFrontendResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fully_qualified_domain_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fully_qualified_domain_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplicationLoadBalancerFrontend_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_frontend", "test")

	r := ApplicationLoadBalancerFrontendResource{}
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

func (r ApplicationLoadBalancerFrontendResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApplicationLoadBalancerFrontendResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
  }
}

%s

resource "azurerm_application_load_balancer_frontend" "test" {
  name                         = "acct-frnt-%d"
  application_load_balancer_id = azurerm_application_load_balancer.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationLoadBalancerFrontendResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
  }
}

%s

resource "azurerm_application_load_balancer_frontend" "test" {
  name                         = "acct-frnt-%d"
  application_load_balancer_id = azurerm_application_load_balancer.test.id
  tags = {
    "tag1" = "value1"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationLoadBalancerFrontendResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_application_load_balancer_frontend" "import" {
  name                         = azurerm_application_load_balancer_frontend.test.name
  application_load_balancer_id = azurerm_application_load_balancer_frontend.test.application_load_balancer_id
}
`, r.basic(data))
}
