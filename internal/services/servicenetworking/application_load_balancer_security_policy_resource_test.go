package servicenetworking_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/securitypoliciesinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApplicationLoadBalancerSecurityPoliciesResource struct{}

func TestAccApplicationLoadBalancerSecurityPolicies_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_security_policy", "test")

	r := ApplicationLoadBalancerSecurityPoliciesResource{}
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

func TestAccApplicationLoadBalancerSecurityPolicies_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_security_policy", "test")

	r := ApplicationLoadBalancerSecurityPoliciesResource{}
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

func TestAccApplicationLoadBalancerSecurityPolicies_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_security_policy", "test")

	r := ApplicationLoadBalancerSecurityPoliciesResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
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

func TestAccApplicationLoadBalancerSecurityPolicies_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_load_balancer_security_policy", "test")

	r := ApplicationLoadBalancerSecurityPoliciesResource{}
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

func (r ApplicationLoadBalancerSecurityPoliciesResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := securitypoliciesinterface.ParseSecurityPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ServiceNetworking.SecurityPoliciesInterface.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ApplicationLoadBalancerSecurityPoliciesResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-alb-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_load_balancer" "test" {
  name                = "acctest-alb-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_web_application_firewall_policy" "test" {
  name                = "acctest-wafpolicy-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  managed_rules {
    managed_rule_set {
      type    = "Microsoft_DefaultRuleSet"
      version = "2.1"
    }
  }

  policy_settings {
    enabled = true
    mode    = "Detection"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApplicationLoadBalancerSecurityPoliciesResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_application_load_balancer_security_policy" "test" {
  name                               = "acctest-albsp-%d"
  application_load_balancer_id       = azurerm_application_load_balancer.test.id
  location                           = azurerm_resource_group.test.location
  web_application_firewall_policy_id = azurerm_web_application_firewall_policy.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationLoadBalancerSecurityPoliciesResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_application_load_balancer_security_policy" "test" {
  name                               = "acctest-albsp-%d"
  application_load_balancer_id       = azurerm_application_load_balancer.test.id
  location                           = azurerm_resource_group.test.location
  web_application_firewall_policy_id = azurerm_web_application_firewall_policy.test.id
  tags = {
    test = "update"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationLoadBalancerSecurityPoliciesResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_application_load_balancer_security_policy" "test" {
  name                               = "acctest-albsp-%d"
  application_load_balancer_id       = azurerm_application_load_balancer.test.id
  location                           = azurerm_resource_group.test.location
  web_application_firewall_policy_id = azurerm_web_application_firewall_policy.test.id
  tags = {
    foo  = "bar"
    test = "complete"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ApplicationLoadBalancerSecurityPoliciesResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_application_load_balancer_security_policy" "import" {
  name                               = azurerm_application_load_balancer_security_policy.test.name
  application_load_balancer_id       = azurerm_application_load_balancer_security_policy.test.application_load_balancer_id
  location                           = azurerm_application_load_balancer_security_policy.test.location
  web_application_firewall_policy_id = azurerm_application_load_balancer_security_policy.test.web_application_firewall_policy_id
}
`, r.basic(data))
}
