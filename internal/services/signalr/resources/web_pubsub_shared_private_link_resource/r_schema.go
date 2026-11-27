// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_shared_private_link_resource

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func webPubsubSharedPrivateLinkResourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Web Pubsub Shared Private Link Resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"web_pubsub_id": commonschema.ResourceIDReferenceRequiredForceNew(&webpubsub.WebPubSubId{}, "Web Pub/Sub Service"),

		"subresource_name": {
			Description:  "The sub resource name which the Web Pubsub Private Endpoint can connect to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PrivateLinkSubResourceName,
		},

		"target_resource_id": {
			Description:  "The ID of the Shared Private Link Enabled Remote Resource which this Web Pubsub Private Endpoint should be connected to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"request_message": {
			Description:  "The request message for requesting approval of the Shared Private Link Enabled Remote Resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"status": {
			Description: "The status of a Private Endpoint connection.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
