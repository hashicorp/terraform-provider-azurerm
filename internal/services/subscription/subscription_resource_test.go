// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package subscription_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubscriptionResource struct{}

func TestAccSubscriptionResource_basic(t *testing.T) {
	if os.Getenv("ARM_BILLING_ACCOUNT") == "" {
		t.Skip("skipping tests - no billing account data provided")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription", "test")
	r := SubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicEnrollmentAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicEnrollmentAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicEnrollmentAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep("billing_scope_id"),
		{
			Config: r.basicEnrollmentAccountUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep("billing_scope_id"),
	})
}

func TestAccSubscriptionResource_devTest(t *testing.T) {
	if os.Getenv("ARM_BILLING_ACCOUNT") == "" {
		t.Skip("skipping tests - no billing account data provided")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription", "test")
	r := SubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicEnrollmentAccountDevTest(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccSubscriptionResource_withOwner(t *testing.T) {
	if os.Getenv("ARM_BILLING_ACCOUNT") == "" {
		t.Skip("skipping tests - no billing account data provided")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription", "test")
	r := SubscriptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicEnrollmentAccountWithOwner(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func (SubscriptionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := subscriptions.ParseAliasID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Subscription.AliasClient.AliasGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Subscription Alias %q: %+v", id.AliasName, err)
	}

	return utils.Bool(true), nil
}

// TODO - Need Env vars in CI for Billing Account and Enrollment Account - Testing disabled for now
func (SubscriptionResource) basicEnrollmentAccount(data acceptance.TestData) string {
	billingAccount := os.Getenv("ARM_BILLING_ACCOUNT")
	billingProfile := os.Getenv("ARM_BILLING_PROFILE")
	invoiceSection := os.Getenv("ARM_INVOICE_SECTION")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_billing_mca_account_scope" "test" {
  billing_account_name = "%[1]s"
  billing_profile_name = "%[2]s"
  invoice_section_name = "%[3]s"
}

resource "azurerm_subscription" "test" {
  alias             = "testAcc-%[4]d"
  subscription_name = "testAccSubscription %[4]d"
  billing_scope_id  = data.azurerm_billing_mca_account_scope.test.id
}
`, billingAccount, billingProfile, invoiceSection, data.RandomInteger)
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

  tags = {
    key = "value"
  }
}
`, billingAccount, enrollmentAccount, data.RandomInteger)
}

func (SubscriptionResource) basicEnrollmentAccountDevTest(data acceptance.TestData) string {
	billingAccount := os.Getenv("ARM_BILLING_ACCOUNT")
	enrollmentAccount := os.Getenv("ARM_BILLING_ENROLLMENT_ACCOUNT")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_billing_enrollment_account_scope" "test" {
  billing_account_name    = "%s"
  enrollment_account_name = "%s"
}

resource "azurerm_subscription" "test" {
  alias             = "testAcc-%[3]d"
  subscription_name = "testAccSubscription Renamed %[3]d"
  billing_scope_id  = data.azurerm_billing_enrollment_account_scope.test.id
  workload          = "DevTest"
}
`, billingAccount, enrollmentAccount, data.RandomInteger)
}

func (SubscriptionResource) basicEnrollmentAccountWithOwner(data acceptance.TestData) string {
	billingAccount := os.Getenv("ARM_BILLING_ACCOUNT")
	enrollmentAccount := os.Getenv("ARM_BILLING_ENROLLMENT_ACCOUNT")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azuread_user" "test" {
  user_principal_name = "test@example.com"
  display_name        = "Test"
}

data "azurerm_billing_enrollment_account_scope" "test" {
  billing_account_name    = "%s"
  enrollment_account_name = "%s"
}

resource "azurerm_subscription" "test" {
  alias                  = "testAcc-%[3]d"
  subscription_name      = "testAccSubscription Renamed %[3]d"
  billing_scope_id       = data.azurerm_billing_enrollment_account_scope.test.id
  subscription_owner_id  = resource.azuread_user.test.object_id
  subscription_tenant_id = data.azurerm_client_config.current.tenant_id
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
