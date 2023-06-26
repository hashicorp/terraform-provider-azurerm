// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

// CheckDestroyedFunc returns a TestCheckFunc which validates the resource no longer exists
func CheckDestroyedFunc(client *clients.Client, testResource types.TestResource, resourceType, resourceName string) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		// even with rate limiting - an exists function should never take more than 5m, so should be safe
		ctx, cancel := context.WithDeadline(client.StopContext, time.Now().Add(5*time.Minute))
		defer cancel()

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
