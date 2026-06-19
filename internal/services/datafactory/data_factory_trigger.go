// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

// startDataFactoryTrigger retries Start because CreateOrUpdate can return success while
// the underlaying Event Grid subscription is still provisioning, causing Start to fail.
func startDataFactoryTrigger(ctx context.Context, client datafactory.TriggersClient, id parse.TriggerId) error {
	for retries := 0; retries < 5; retries++ {
		future, err := client.Start(ctx, id.ResourceGroup, id.FactoryName, id.Name)
		if err == nil {
			return future.WaitForCompletionRef(ctx, client.Client)
		}

		if !strings.Contains(err.Error(), "Resource cannot be updated during provisioning") {
			return err
		}

		if retries < 4 {
			time.Sleep(5 * time.Second)
		}
	}

	return fmt.Errorf("starting %s: trigger remained in provisioning", id)
}
