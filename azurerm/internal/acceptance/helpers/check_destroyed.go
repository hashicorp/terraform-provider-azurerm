package helpers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

// CheckDestroyedFunc returns a TestCheckFunc which validates the resource no longer exists
func CheckDestroyedFunc(client *clients.Client, testResource types.TestResource, resourceType, resourceName string) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		ctx := client.StopContext

		for label, resourceState := range state.RootModule().Resources {
			if resourceState.Type != resourceType {
				continue
			}
			if label != resourceName {
				continue
			}

			// Destroy is unconcerned with an error checking the status, since this is going to be "not found"
			result, err := testResource.Exists(ctx, client, resourceState.Primary)
			if result == nil && err == nil {
				return fmt.Errorf("should have either an error or a result when checking if %q has been destroyed", resourceName)
			}
			if result != nil && *result {
				return fmt.Errorf("%q still exists", resourceName)
			}
		}

		return nil
	}
}
