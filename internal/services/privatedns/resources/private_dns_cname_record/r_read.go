// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_cname_record

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateDnsCNameRecordRead(d *pluginsdk.ResourceData, meta any) error {
	dnsClient := meta.(*clients.Client).PrivateDns.RecordSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatedns.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := dnsClient.RecordSetsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return resourcePrivateDnsCNameRecordFlatten(d, id, resp.Model)
}

func resourcePrivateDnsCNameRecordImporter(_ context.Context, d *pluginsdk.ResourceData, _ any) ([]*pluginsdk.ResourceData, error) {
	resourceId, err := privatedns.ParseRecordTypeID(d.Id())
	if err != nil {
		return []*pluginsdk.ResourceData{d}, err
	}
	if resourceId.RecordType != privatedns.RecordTypeCNAME {
		return []*pluginsdk.ResourceData{d}, fmt.Errorf("importing %s wrong type received: expected %s received %s", resourceId, privatedns.RecordTypeCNAME, resourceId.RecordType)
	}
	return []*pluginsdk.ResourceData{d}, nil
}
