// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_object_replication

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageObjectReplicationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.ResourceManager.ObjectReplicationPolicyOperationGroup
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ObjectReplicationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.ObjectReplicationPoliciesDelete(ctx, id.Dst); err != nil {
		return fmt.Errorf("deleting %q: %+v", id.Dst, err)
	}

	if _, err := client.ObjectReplicationPoliciesDelete(ctx, id.Src); err != nil {
		return fmt.Errorf("deleting %q : %+v", id.Dst, err)
	}
	return nil
}
