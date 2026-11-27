// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_protected_file_share

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	recoveryservicesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectedFileShareDelete(d *pluginsdk.ResourceData, meta any) error {
	client := recoveryservicesClient.From(meta).ProtectedItemsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
