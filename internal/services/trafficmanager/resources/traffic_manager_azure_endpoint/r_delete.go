// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_azure_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAzureEndpointDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).TrafficManager.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
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
