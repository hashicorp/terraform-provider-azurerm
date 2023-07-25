// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

func resourceIotHubEnrichment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmIotHubEnrichmentCreateUpdate,
		Read:   resourceArmIotHubEnrichmentRead,
		Update: resourceArmIotHubEnrichmentCreateUpdate,
		Delete: resourceArmIotHubEnrichmentDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.IoTHubEnrichmentV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EnrichmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"key": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
					"Enrichment Key name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"iothub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"value": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"endpoint_names": {
				Type:     pluginsdk.TypeList,
				MaxItems: 100,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Required: true,
			},
		},
	}
}

func resourceArmIotHubEnrichmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	enrichmentKey := d.Get("key").(string)
	enrichmentValue := d.Get("value").(string)
	endpointNamesRaw := d.Get("endpoint_names").([]interface{})

	enrichment := devices.EnrichmentProperties{
		Key:           &enrichmentKey,
		Value:         &enrichmentValue,
		EndpointNames: utils.ExpandStringSlice(endpointNamesRaw),
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
					return tf.ImportAsExistsError("azurerm_iothub_enrichment", id.ID())
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
		return fmt.Errorf("Unable to find Enrichment %q defined for IotHub %q (Resource Group %q)", enrichmentKey, iothubName, resourceGroup)
	}
	routing.Enrichments = &enrichments

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("creating/updating IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceArmIotHubEnrichmentRead(d, meta)
}

func resourceArmIotHubEnrichmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EnrichmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
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

func resourceArmIotHubEnrichmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EnrichmentID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}
		return fmt.Errorf("retrieving IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	if iothub.Properties == nil || iothub.Properties.Routing == nil {
		return nil
	}

	enrichments := iothub.Properties.Routing.Enrichments
	if enrichments == nil {
		return nil
	}

	updatedEnrichments := make([]devices.EnrichmentProperties, 0)
	for _, enrichment := range *enrichments {
		if enrichment.Key != nil {
			if !strings.EqualFold(*enrichment.Key, id.Name) {
				updatedEnrichments = append(updatedEnrichments, enrichment)
			}
		}
	}
	iothub.Properties.Routing.Enrichments = &updatedEnrichments

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
