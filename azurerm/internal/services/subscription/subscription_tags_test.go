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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SubscriptionTags struct{}

func TestSubscriptionTags_basic(t *testing.T) {
	if os.Getenv("ARM_SUBSCRIPTION_ID") == "" {
		t.Skip("skipping tests - no subscription ID data provided")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription_tags", "test")
	r := SubscriptionTags{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestSubscriptionTags_requiresImport(t *testing.T) {
	if os.Getenv("ARM_SUBSCRIPTION_ID") == "" {
		t.Skip("skipping tests - no subscription ID data provided")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription_tags", "test")
	r := SubscriptionTags{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.RequiresImportErrorStep(r.requiresImportConfig),
	})
}

func TestSubscriptionTags_updateWithTags(t *testing.T) {
	if os.Getenv("ARM_SUBSCRIPTION_ID") == "" {
		t.Skip("skipping tests - no subscription ID data provided")
	}

	data := acceptance.BuildTestData(t, "azurerm_subscription_tags", "test")
	r := SubscriptionTags{}
	assert := check.That(data.ResourceName)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicConfig(data),
			Check: resource.ComposeTestCheckFunc(
				assert.ExistsInAzure(r),
				assert.Key("tags.%").HasValue("2"),
				assert.Key("tags.cost_center").HasValue("MSFT"),
				assert.Key("tags.environment").HasValue("Production"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTagsUpdatedConfig(data),
			Check: resource.ComposeTestCheckFunc(
				assert.ExistsInAzure(r),
				assert.Key("tags.%").HasValue("1"),
				assert.Key("tags.environment").HasValue("staging"),
			),
		},
		data.ImportStep(),
	})
}

func (t SubscriptionTags) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	subscriptionId := state.Attributes["subscription_id"]
	atags := state.Attributes["tags"]
	fmt.Printf("tags: %T\n %s\n", atags, atags)
	fmt.Println("subscription id", subscriptionId, "tags", atags, "\n ")
	resp, err := client.Subscription.TagsClient.GetAtScope(ctx, "subscriptions/"+subscriptionId)
	if err != nil {
		return nil, fmt.Errorf("retrieving tags from subscription %q: %+v", subscriptionId, err)
	}
	fmt.Println(resp.Properties)
	return utils.Bool(resp.Properties != nil), nil
}

func (t SubscriptionTags) basicConfig(data acceptance.TestData) string {
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_subscription_tags" "test" {
  subscription_id = "%s"
  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
		`, subscriptionId)
}

func (t SubscriptionTags) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
resource "azurerm_subscription_tags" "import" {
  subscription_id = azurerm_subscription_tags.test.subscription_id
  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
	`, t.basicConfig(data))
}

func (t SubscriptionTags) withTagsUpdatedConfig(data acceptance.TestData) string {
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_tags" "test" {
  subscription_id = "%s"
  tags = {
    environment = "staging"
  }
}
	`, subscriptionId)
}
