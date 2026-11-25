// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_endpoint_eventgrid

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func digitalTwinsEndpointEventgridSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Digital Twins Eventgrid Endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DigitalTwinsInstanceName,
		},

		"digital_twins_id": {
			Description:  "The resource ID of the Digital Twins Instance.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: digitaltwinsinstance.ValidateDigitalTwinsInstanceID,
		},

		"eventgrid_topic_endpoint": {
			Description:  "The endpoint of the Event Grid Topic.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"eventgrid_topic_primary_access_key": {
			Description:  "The primary access key of the Event Grid Topic.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"eventgrid_topic_secondary_access_key": {
			Description:  "The secondary access key of the Event Grid Topic.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dead_letter_storage_secret": {
			Description:  "The storage secret of the dead-lettering, whose format is `https://<storageAccountname>.blob.core.windows.net/<containerName>?<SASToken>`. When an endpoint can't deliver an event within a certain time period or after trying to deliver the event a certain number of times, it can send the undelivered event to a Storage Account.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
