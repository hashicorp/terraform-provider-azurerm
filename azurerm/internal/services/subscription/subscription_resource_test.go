package subscription_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SubscriptionResource struct{}

func TestAccSubscriptionResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription", "test")
	r := SubscriptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicEnrollmentAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
	})
}

func TestAccSubscriptionResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription", "test")
	r := SubscriptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.requiresImport(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
	})
}

func (SubscriptionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SubscriptionAliasID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Subscription.AliasClient.Get(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Subscription Alias %q: %+v", id.Name, err)
	}

	return utils.Bool(true), nil
}

// TODO - Need Env vars for Billing Account and Enrollment Account
func (SubscriptionResource) basicEnrollmentAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription" "test" {
  alias              = "testAcc-%[1]d"
  subscription_name  = "testAccSubscription %[1]d"
  billing_account    = ""
  enrollment_account = ""
}
`, data.RandomInteger, data.RandomString)
}
func (r SubscriptionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subscription" "import" {
  alias              = azurerm_subscription.test.alias
  subscription_name  = azurerm_subscription.test.subscription_name
  billing_account    = azurerm_subscription.test.billing_account
  enrollment_account = azurerm_subscription.test.enrollment_account
}
`, r.basicEnrollmentAccount(data))
}
