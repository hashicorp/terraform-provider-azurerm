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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
			// even with rate limiting - an exists function should never take more than 5m, so should be safe
			ctx, cancel := context.WithDeadline(client.StopContext, time.Now().Add(5*time.Minute))
			defer cancel()

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
