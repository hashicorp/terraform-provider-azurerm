// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_ptr_record

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateDnsPtrRecordDelete(d *pluginsdk.ResourceData, meta any) error {
	dnsClient := client.From(meta).RecordSetsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := privatedns.ParseRecordTypeID(d.Id())
	if err != nil {
		return err
	}

	options := privatedns.RecordSetsDeleteOperationOptions{IfMatch: pointer.To("")}

	if _, err = dnsClient.RecordSetsDelete(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
