// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custom_provider

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	customprovidersClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/customproviders/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCustomProviderDelete(d *pluginsdk.ResourceData, meta any) error {
	client := customprovidersClient.From(meta).CustomProviderClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := customresourceprovider.ParseResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
