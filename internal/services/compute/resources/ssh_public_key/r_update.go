// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ssh_public_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/sshpublickeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSshPublicKeyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.SSHPublicKeysClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := sshpublickeys.ParseSshPublicKeyID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Get(ctx, *id); err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	payload := sshpublickeys.SshPublicKeyUpdateResource{}
	if d.HasChange("public_key") {
		payload.Properties = &sshpublickeys.SshPublicKeyResourceProperties{
			PublicKey: pointer.To(d.Get("public_key").(string)),
		}
	}
	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]any)
		payload.Tags = tags.Expand(tagsRaw)
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceSshPublicKeyRead(d, meta)
}
