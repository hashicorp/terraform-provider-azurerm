// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_custom_https_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	frontdoorClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorCustomHTTPSConfigurationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := frontdoorClient.From(meta).FrontDoorsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := frontdoors.ParseFrontendEndpointIDInsensitively(d.Get("frontend_endpoint_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.FrontendEndpointsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if model.Properties == nil {
			return fmt.Errorf("reading %s: `properties` was nil", id)
		}
		props := *model.Properties

		input := customHttpsConfigurationUpdateInput{
			customHttpsConfigurationCurrent: props.CustomHTTPSConfiguration,
			customHttpsConfigurationNew:     make([]any, 0),
			frontendEndpointId:              *id,
		}

		if props.CustomHTTPSProvisioningState != nil {
			input.provisioningState = *props.CustomHTTPSProvisioningState
		}
		if err := updateCustomHTTPSConfiguration(ctx, client, input); err != nil {
			return fmt.Errorf("disabling Custom HTTPS configuration for Frontend Endpoint %q (Front Door %q / Resource Group %q): %+v", id.FrontendEndpointName, id.FrontDoorName, id.ResourceGroupName, err)
		}
	}

	return nil
}
