package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRuleStackResource struct{}

func TestAccPaloAltoLocalRuleStack_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rule_stack", "test")

	r := LocalRuleStackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{{
		Config: r.basic(data),
		Check:  acceptance.ComposeTestCheckFunc(),
	}})
}

func (r LocalRuleStackResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := localrulestacks.ParseLocalRuleStackID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.LocalRuleStacksClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r LocalRuleStackResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rule_stack" "test" {
  name                = "testAcc-palrs-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
}

`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r LocalRuleStackResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-CAE-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
