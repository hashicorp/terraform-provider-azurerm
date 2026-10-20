// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_share_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
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
