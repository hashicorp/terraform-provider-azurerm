package paloalto_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRuleResource struct{}

func TestAccPaloAltoLocalRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack_rule", "test")

	r := LocalRuleResource{}

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

func TestAccPaloAltoLocalRule_withDestination(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack_rule", "test")

	r := LocalRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDestination(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPaloAltoLocalRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack_rule", "test")

	r := LocalRuleResource{}

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

func TestAccPaloAltoLocalRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack_rule", "test")

	r := LocalRuleResource{}

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

func (r LocalRuleResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := localrules.ParseLocalRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.LocalRulesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r LocalRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rule_stack_rule" "test" {
  name          = "testacc-palr-%[2]d"
  rule_stack_id = azurerm_palo_alto_local_rule_stack.test.id
  priority      = 100

  applications = ["any"]

  destination {
    cidrs = ["any"]
  }

  source {
    cidrs = ["any"]
  }
}


`, r.template(data), data.RandomInteger)
}

func (r LocalRuleResource) withDestination(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rule_stack_rule" "test" {
  name          = "testacc-palr-%[2]d"
  rule_stack_id = azurerm_palo_alto_local_rule_stack.test.id
  priority      = 100

  applications = ["any"]

  destination {
	countries = ["US", "GB"]
  }

  source {
	countries = ["US", "GB"]
  }
}


`, r.template(data), data.RandomInteger)
}

func (r LocalRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rule_stack_certificate" "test" {
  name          = "testacc-palc-%[2]d"
  rule_stack_id = azurerm_palo_alto_local_rule_stack.test.id
  self_signed   = true
}


resource "azurerm_palo_alto_local_rule_stack_rule" "test" {
  name          = "testacc-palr-%[2]d"
  rule_stack_id = azurerm_palo_alto_local_rule_stack.test.id
  priority      = 100

  action = "DenySilent"
  applications = ["any"]
  audit_comment = "test audit comment"

  //category {
  //  // feeds = ["foo", "bar"] // Needs feeds defined on the LocalRulestack?
  //  // custom_urls = ["https://microsoft.com"] // TODO - This is another resource type in PAN?
  //}

  decryption_rule_type = "SSLOutboundInspection" // TODO - Needs Certs to be available on the RuleStack
  description = "Acceptance Test Rule - dated %[2]d"

  destination {
	countries = ["US", "GB"]
  }

  logging_enabled = false

  inspection_certificate_id = azurerm_palo_alto_local_rule_stack_certificate.test.id

  negate_destination = true
  negate_source = true
  
  protocol = "TCP:8080"
  
  enabled = false

  source {
	countries = ["US", "GB"]
  }

  tags = {
    "acctest" = "true"
    "foo" = "bar"
  }
}


`, r.template(data), data.RandomInteger)
}

func (r LocalRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PAN-%[1]d"
  location = "%[2]s"
}

resource "azurerm_palo_alto_local_rule_stack" "test" {
  name                = "testAcc-palrs-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
