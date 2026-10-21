// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor_custom_https_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorCustomHTTPSConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := frontdoors.ParseFrontendEndpointIDInsensitively(d.Get("frontend_endpoint_id").(string))
	if err != nil {
		return err
	}

	customHttpsConfigurationId := parse.NewCustomHttpsConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.FrontDoorName, id.FrontendEndpointName)

	resp, err := client.FrontendEndpointsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return fmt.Errorf("reading %s: `properties` was nil", id)
	}
	props := *resp.Model.Properties

	input := customHttpsConfigurationUpdateInput{
		customHttpsConfigurationCurrent: props.CustomHTTPSConfiguration,
		customHttpsConfigurationNew:     d.Get("custom_https_configuration").([]any),
		customHttpsProvisioningEnabled:  d.Get("custom_https_provisioning_enabled").(bool),
		frontendEndpointId:              *id,
	}

	if props.CustomHTTPSProvisioningState != nil {
		input.provisioningState = *props.CustomHTTPSProvisioningState
	}

	if err := updateCustomHTTPSConfiguration(ctx, client, input); err != nil {
		return fmt.Errorf("updating Custom HTTPS configuration for %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(customHttpsConfigurationId.ID())
	}

	return resourceFrontDoorCustomHTTPSConfigurationRead(d, meta)
}
