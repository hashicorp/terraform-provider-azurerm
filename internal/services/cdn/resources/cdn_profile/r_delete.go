// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_profile

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	cdnClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnProfileDelete(d *pluginsdk.ResourceData, meta any) error {
	client := cdnClient.From(meta).ProfilesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := profiles.ParseProfileID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.ProfileName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}
