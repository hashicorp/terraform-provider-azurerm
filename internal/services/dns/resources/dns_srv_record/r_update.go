// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_srv_record

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDnsSrvRecordUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Dns.RecordSets
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := recordsets.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := *existing.Model

	if d.HasChange("record") {
		payload.Properties.SRVRecords = expandAzureRmDnsSrvRecords(d)
	}

	if d.HasChange("ttl") {
		payload.Properties.TTL = pointer.To(int64(d.Get("ttl").(int)))
	}

	if d.HasChange("tags") {
		payload.Properties.Metadata = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, payload, recordsets.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDnsSrvRecordRead(d, meta)
}
