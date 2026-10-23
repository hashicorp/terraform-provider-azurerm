// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_enrichment

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceArmIotHubEnrichmentRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EnrichmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[DEBUG] IoTHub %q was not found in Resource Group %q (so Enrichment cannot exist) - removing from state", id.IotHubName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	var props *devices.EnrichmentProperties
	if resp.Properties != nil && resp.Properties.Routing != nil && resp.Properties.Routing.Enrichments != nil {
		for _, enrichment := range *resp.Properties.Routing.Enrichments {
			if enrichment.Key != nil {
				if strings.EqualFold(*enrichment.Key, id.Name) {
					props = &enrichment
					break
				}
			}
		}
	}

	if props == nil {
		log.Printf("[DEBUG] %s was not found - removing from state", *id)
		d.SetId("")
		return nil
	}

	d.Set("key", id.Name)
	d.Set("iothub_name", id.IotHubName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("value", props.Value)
	d.Set("endpoint_names", props.EndpointNames)

	return nil
}
