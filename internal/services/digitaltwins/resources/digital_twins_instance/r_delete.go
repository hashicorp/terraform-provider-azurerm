// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_instance

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	digitaltwinsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDigitalTwinsInstanceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := digitaltwinsClient.From(meta).InstanceClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := digitaltwinsinstance.ParseDigitalTwinsInstanceID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DigitalTwinsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
