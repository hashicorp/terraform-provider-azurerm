// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_nested_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	trafficmanagerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNestedEndpointDelete(d *pluginsdk.ResourceData, meta any) error {
	client := trafficmanagerClient.From(meta).EndpointsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := trafficmanagers.ParseEndpointTypeID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.EndpointsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
