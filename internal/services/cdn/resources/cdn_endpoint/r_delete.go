// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	cdnClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnEndpointDelete(d *pluginsdk.ResourceData, meta any) error {
	client := cdnClient.From(meta).EndpointsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := endpoints.ParseEndpointID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.ProfileName, id.EndpointName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
