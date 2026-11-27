// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afdendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	cdnClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorEndpointDelete(d *pluginsdk.ResourceData, meta any) error {
	client := cdnClient.From(meta).AFDEndpointsClient

	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
