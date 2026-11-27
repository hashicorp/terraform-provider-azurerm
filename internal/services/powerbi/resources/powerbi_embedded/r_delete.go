// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package powerbi_embedded

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	powerbiClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/powerbi/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePowerBIEmbeddedDelete(d *pluginsdk.ResourceData, meta any) error {
	client := powerbiClient.From(meta).CapacityClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := capacities.ParseCapacityID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
