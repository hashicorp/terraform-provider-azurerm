// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/capacitypools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppPoolUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacitypools.ParseCapacityPoolID(d.Id())
	if err != nil {
		return err
	}

	update := capacitypools.CapacityPoolPatch{
		Properties: &capacitypools.PoolPatchProperties{},
	}

	if d.HasChange("size_in_tb") {
		sizeInTB := int64(d.Get("size_in_tb").(int))
		sizeInMB := sizeInTB * 1024 * 1024
		sizeInBytes := sizeInMB * 1024 * 1024

		update.Properties.Size = pointer.To(sizeInBytes)
	}

	if d.HasChange("qos_type") {
		qosType := capacitypools.QosType(d.Get("qos_type").(string))
		update.Properties.QosType = &qosType
	}

	if d.HasChange("cool_access_enabled") {
		update.Properties.CoolAccess = pointer.To(d.Get("cool_access_enabled").(bool))
	}

	if d.HasChange("custom_throughput_mibps") {
		if customThroughputMibps, ok := d.GetOk("custom_throughput_mibps"); ok {
			update.Properties.CustomThroughputMibps = pointer.To(int64(customThroughputMibps.(int)))
		}
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]any)
		update.Tags = tags.Expand(tagsRaw)
	}

	if err = client.PoolsUpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", id.ID(), err)
	}

	return resourceNetAppPoolRead(d, meta)
}
