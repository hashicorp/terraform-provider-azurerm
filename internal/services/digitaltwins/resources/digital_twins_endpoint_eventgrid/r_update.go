// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_endpoint_eventgrid

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDigitalTwinsEndpointEventGridCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.EndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	digitalTwinsId, err := digitaltwinsinstance.ParseDigitalTwinsInstanceID(d.Get("digital_twins_id").(string))
	if err != nil {
		return err
	}

	id := endpoints.NewEndpointID(subscriptionId, digitalTwinsId.ResourceGroupName, digitalTwinsId.DigitalTwinsInstanceName, d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.DigitalTwinsEndpointGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_digital_twins_endpoint_eventgrid", id.ID())
			}
		}
	}

	payload := endpoints.DigitalTwinsEndpointResource{
		Properties: &endpoints.EventGrid{
			AuthenticationType: pointer.To(endpoints.AuthenticationTypeKeyBased),
			TopicEndpoint:      d.Get("eventgrid_topic_endpoint").(string),
			AccessKey1:         d.Get("eventgrid_topic_primary_access_key").(string),
			AccessKey2:         pointer.To(d.Get("eventgrid_topic_secondary_access_key").(string)),
			DeadLetterSecret:   pointer.To(d.Get("dead_letter_storage_secret").(string)),
		},
	}

	if d.IsNewResource() {
		if err := client.DigitalTwinsEndpointCreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.DigitalTwinsEndpointCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}
	return resourceDigitalTwinsEndpointEventGridRead(d, meta)
}
