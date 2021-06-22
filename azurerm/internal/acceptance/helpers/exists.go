package helpers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func DoesNotExistInAzure(client *clients.Client, testResource types.TestResource, resourceName string) pluginsdk.TestCheckFunc {
	return existsFunc(false)(client, testResource, resourceName)
}

func ExistsInAzure(client *clients.Client, testResource types.TestResource, resourceName string) pluginsdk.TestCheckFunc {
	return existsFunc(true)(client, testResource, resourceName)
}

func existsFunc(shouldExist bool) func(*clients.Client, types.TestResource, string) pluginsdk.TestCheckFunc {
	return func(client *clients.Client, testResource types.TestResource, resourceName string) pluginsdk.TestCheckFunc {
		return func(s *terraform.State) error {
			ctx := client.StopContext

			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("%q was not found in the state", resourceName)
			}

			result, err := testResource.Exists(ctx, client, rs.Primary)
			if err != nil {
				return fmt.Errorf("running exists func for %q: %+v", resourceName, err)
			}
			if result == nil {
				return fmt.Errorf("received nil for exists for %q", resourceName)
			}

			if *result != shouldExist {
				if !shouldExist {
					return fmt.Errorf("%q still exists", resourceName)
				}

				return fmt.Errorf("%q did not exist", resourceName)
			}

			return nil
		}
	}
}
