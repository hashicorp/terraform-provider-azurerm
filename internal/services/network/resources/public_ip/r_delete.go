// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package public_ip

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePublicIpDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).PublicIPAddresses
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParsePublicIPAddressID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
