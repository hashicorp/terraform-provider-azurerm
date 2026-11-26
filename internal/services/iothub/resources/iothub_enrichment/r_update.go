// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_enrichment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceArmIotHubEnrichmentCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	iothubName := d.Get("iothub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(iothubName, IothubResourceName)
	defer locks.UnlockByName(iothubName, IothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if response.WasNotFound(iothub.Response.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	enrichmentKey := d.Get("key").(string)
	enrichmentValue := d.Get("value").(string)
	endpointNamesRaw := d.Get("endpoint_names").([]any)

	enrichment := devices.EnrichmentProperties{
		Key:           &enrichmentKey,
		Value:         &enrichmentValue,
		EndpointNames: pluginsdk.ExpandStringSlice(endpointNamesRaw),
	}

	routing := iothub.Properties.Routing
	if routing == nil {
		routing = &devices.RoutingProperties{}
	}

	if routing.Enrichments == nil {
		enrichments := make([]devices.EnrichmentProperties, 0)
		routing.Enrichments = &enrichments
	}

	enrichments := make([]devices.EnrichmentProperties, 0)

	id := parse.NewEnrichmentID(subscriptionId, resourceGroup, iothubName, enrichmentKey)
	alreadyExists := false
	for _, existingEnrichment := range *routing.Enrichments {
		if existingEnrichment.Key != nil {
			if strings.EqualFold(*existingEnrichment.Key, enrichmentKey) {
				if d.IsNewResource() {
					if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError("azurerm_iothub_enrichment", id.ID())
					}
				}
				enrichments = append(enrichments, enrichment)
				alreadyExists = true
			} else {
				enrichments = append(enrichments, existingEnrichment)
			}
		}
	}

	if d.IsNewResource() {
		enrichments = append(enrichments, enrichment)
	} else if !alreadyExists {
		return fmt.Errorf("unable to find Enrichment %q defined for IotHub %q (Resource Group %q)", enrichmentKey, iothubName, resourceGroup)
	}
	routing.Enrichments = &enrichments

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("creating/updating IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	return resourceArmIotHubEnrichmentRead(d, meta)
}
