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

// DeleteResourceFunc returns a TestCheckFunc which deletes the resource within Azure
// this is only used within the Internal
func DeleteResourceFunc(client *clients.Client, testResource types.TestResourceVerifyingRemoved, resourceName string) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		// @tombuildsstuff: the delete function shouldn't take more than an hour
		// on the off-chance that it's not for a given resource, we may want to add an optional interface
		// to return the deletion timeout, rather than extending this for all resources
		ctx, cancel := context.WithDeadline(client.StopContext, time.Now().Add(1*time.Hour))
		defer cancel()

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%q was not found in the state", resourceName)
		}

		result, err := testResource.Destroy(ctx, client, rs.Primary)
		if err != nil {
			return fmt.Errorf("running destroy func for %q: %+v", resourceName, err)
		}
		if result == nil {
			return fmt.Errorf("received nil for destroy result for %q", resourceName)
		}

		if !*result {
			return fmt.Errorf("deleting %q but no error", resourceName)
		}

		return nil
	}
}
