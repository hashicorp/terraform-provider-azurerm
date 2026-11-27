// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_protection_container

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	recoveryservicesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryProtectionContainerDelete(d *pluginsdk.ResourceData, meta any) error {
	id, err := replicationprotectioncontainers.ParseReplicationProtectionContainerID(d.Id())
	if err != nil {
		return err
	}

	client := recoveryservicesClient.From(meta).ProtectionContainerClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting site recovery protection container %s : %+v", id.String(), err)
	}

	return nil
}
