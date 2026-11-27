// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_disk_sas_token

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedDiskSasTokenDelete(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).DisksClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	if err = client.RevokeAccessThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("revoking access to %s: %+v", *id, err)
	}

	return nil
}
