package subscription_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubscriptionNameResource struct{}

func TestAccSubscriptionNameResource(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests because we are renaming 1 subscription
	// If multiple tests are run in parallel, the test cases will fail.
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"all": {
			"basic":  testAccSubscriptionNameResource_basic,
			"import": testAccSubscriptionNameResource_requiresImport,
			"update": testAccSubscriptionNameResource_update,
		},
	})
}

func testAccSubscriptionNameResource_basic(t *testing.T) {
	if os.Getenv("ARM_RENAME_SUBSCRIPTION_ID") == "" {
		t.Skip("skipping tests - no subscription to rename.")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription_name", "test")
	r := SubscriptionNameResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func testAccSubscriptionNameResource_requiresImport(t *testing.T) {
	if os.Getenv("ARM_RENAME_SUBSCRIPTION_ID") == "" {
		t.Skip("skipping tests - no subscription to rename")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription_name", "test")
	r := SubscriptionNameResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccSubscriptionNameResource_update(t *testing.T) {
	if os.Getenv("ARM_RENAME_SUBSCRIPTION_ID") == "" {
		t.Skip("skipping tests - no subscription to rename")
	}
	data := acceptance.BuildTestData(t, "azurerm_subscription_name", "test")
	r := SubscriptionNameResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.tag1").HasValue("2"),
				check.That(data.ResourceName).Key("tags.newTag").HasValue("hello"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func (SubscriptionNameResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (SubscriptionNameResource) basic(data acceptance.TestData) string {
	renameSubscriptionId := os.Getenv("ARM_RENAME_SUBSCRIPTION_ID")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_name" "test" {
  subscription_id = "%s"
  name            = "subscription-basic-%d"
}
`, renameSubscriptionId, data.RandomInteger)
}

func (SubscriptionNameResource) update(data acceptance.TestData) string {
	renameSubscriptionId := os.Getenv("ARM_RENAME_SUBSCRIPTION_ID")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_name" "test" {
  subscription_id = "%s"
  name            = "subscription-updated-%d"

  tags = {
    tag1   = "2"
    newTag = "hello"
  }
}
`, renameSubscriptionId, data.RandomInteger)
}

func (r SubscriptionNameResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subscription_name" "import" {
  subscription_id = azurerm_subscription_name.test.subscription_id
  name            = azurerm_subscription_name.test.name
}
`, r.basic(data))
}
