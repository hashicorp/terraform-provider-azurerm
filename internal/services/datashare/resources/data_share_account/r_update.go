// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	datashareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := datashareClient.From(meta).AccountClient
	ctx, cancel := timeouts.ForUpdate(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	props := account.AccountUpdateParameters{}

	if d.HasChange("tags") {
		props.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err = client.Update(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDataShareAccountRead(d, meta)
}
