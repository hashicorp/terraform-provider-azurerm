// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_account

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := batchaccount.ParseBatchAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			log.Printf("[DEBUG] Batch Account %s - removing from state!", *id)
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	return resourceBatchAccountFlatten(ctx, client, d, id, resp.Model, true)
}

func isShardKeyAllowed(input []any) bool {
	if len(input) == 0 {
		return false
	}
	for _, authMod := range input {
		if strings.EqualFold(authMod.(string), string(batchaccount.AuthenticationModeSharedKey)) {
			return true
		}
	}
	return false
}
