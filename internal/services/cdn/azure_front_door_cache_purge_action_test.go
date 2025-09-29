package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type AzureFrontDoorCachePurgeAction struct{}

func TestAccAzureFrontDoorCachePurgeAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_front_door_cache_purge", "test")
	a := AzureFrontDoorCachePurgeAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.pathsOnly(data),
				Check:  nil, // TODO - plugin-testing release?
			},
		},
	})
}

func TestAccAzureFrontDoorCachePurgeAction_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azure_front_door_cache_purge", "test")
	a := AzureFrontDoorCachePurgeAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.complete(data),
				Check:  nil, // TODO - plugin-testing release?
			},
		},
	})
}

func (a *AzureFrontDoorCachePurgeAction) pathsOnly(_ acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "terraform_data" "trigger" {
  input = "trigger"
  lifecycle {
    action_trigger {
      events  = [before_create, before_update]
      actions = [action.azurerm_azure_front_door_cache_purge.test]
    }
  }
}

action "azurerm_azure_front_door_cache_purge" "test" {
  config {
    front_door_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1"
    content_paths = [
      "/*"
    ]
  }
}
`)
}

func (a *AzureFrontDoorCachePurgeAction) complete(_ acceptance.TestData) string {
	return fmt.Sprintf(`

provider "azurerm" {
  features {}
}

resource "terraform_data" "trigger" {
  input = "trigger"
  lifecycle {
    action_trigger {
      events  = [before_create, before_update]
      actions = [action.azurerm_azure_front_door_cache_purge.test]
    }
  }
}

action "azurerm_azure_front_door_cache_purge" "test" {
  config {
    front_door_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1"
    content_paths = [
      "/*"
    ]
    domains = [
      "contoso.com"
    ]
  }
}
`)
}
