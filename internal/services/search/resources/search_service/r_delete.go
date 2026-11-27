// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	searchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSearchServiceDelete(d *pluginsdk.ResourceData, meta any) error {
	client := searchClient.From(meta).ServicesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := services.ParseSearchServiceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id, services.DeleteOperationOptions{}); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
