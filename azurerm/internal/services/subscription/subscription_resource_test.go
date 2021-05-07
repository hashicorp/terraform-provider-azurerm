package subscription_test

import (
	"context"
	"fmt"
	"os"
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
	if os.Getenv("ARM_BILLING_ACCOUNT") == "" {
		t.Skip("skipping tests - no billing account data provided")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription", "test")
	r := SubscriptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicEnrollmentAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionResource_requiresImport(t *testing.T) {
	if os.Getenv("ARM_BILLING_ACCOUNT") == "" {
		t.Skip("skipping tests - no billing account data provided")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription", "test")
	r := SubscriptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicEnrollmentAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSubscriptionResource_update(t *testing.T) {
	if os.Getenv("ARM_BILLING_ACCOUNT") == "" {
		t.Skip("skipping tests - no billing account data provided")
	}
	data := acceptance.BuildTestData(t, "azurerm_subscription", "test")
	r := SubscriptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicEnrollmentAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep("billing_scope_id"),
		{
			Config: r.basicEnrollmentAccountUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep("billing_scope_id"),
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

// TODO - Need Env vars in CI for Billing Account and Enrollment Account - Testing disabled for now
func (SubscriptionResource) basicEnrollmentAccount(data acceptance.TestData) string {
	billingAccount := os.Getenv("ARM_BILLING_ACCOUNT")
	enrollmentAccount := os.Getenv("ARM_BILLING_ENROLLMENT_ACCOUNT")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_billing_enrollment_account_scope" "test" {
  billing_account    = "%s"
  enrollment_account = "%s"
}

resource "azurerm_subscription" "test" {
  alias             = "testAcc-%[3]d"
  subscription_name = "testAccSubscription %[3]d"
  billing_scope_id  = data.azurerm_billing_enrollment_account_scope.test.id
}
`, billingAccount, enrollmentAccount, data.RandomInteger)
}

func (SubscriptionResource) basicEnrollmentAccountUpdate(data acceptance.TestData) string {
	billingAccount := os.Getenv("ARM_BILLING_ACCOUNT")
	enrollmentAccount := os.Getenv("ARM_BILLING_ENROLLMENT_ACCOUNT")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_billing_enrollment_account_scope" "test" {
  billing_account    = "%s"
  enrollment_account = "%s"
}

resource "azurerm_subscription" "test" {
  alias             = "testAcc-%[3]d"
  subscription_name = "testAccSubscription Renamed %[3]d"
  billing_scope_id  = data.azurerm_billing_enrollment_account_scope.test.id
}
`, billingAccount, enrollmentAccount, data.RandomInteger)
}

func (r SubscriptionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subscription" "import" {
  alias             = azurerm_subscription.test.alias
  subscription_name = azurerm_subscription.test.subscription_name
  billing_scope_id  = azurerm_subscription.test.billing_scope_id
}
`, r.basicEnrollmentAccount(data))
}
