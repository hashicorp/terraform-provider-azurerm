// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ssh_public_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/sshpublickeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSshPublicKeyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).SSHPublicKeysClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := sshpublickeys.ParseSshPublicKeyID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
