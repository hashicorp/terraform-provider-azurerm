// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_shared_private_link_resource

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages the Shared Private Link Resource for a Web Pubsub service.",
		Create:      resourceWebPubsubSharedPrivateLinkServiceCreateUpdate,
		Read:        resourceWebPubsubSharedPrivateLinkServiceRead,
		Update:      resourceWebPubsubSharedPrivateLinkServiceCreateUpdate,
		Delete:      resourceWebPubsubSharedPrivateLinkServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webpubsub.ParseSharedPrivateLinkResourceID(id)
			return err
		}),

		Schema: webPubsubSharedPrivateLinkResourceSchema(),
	}
}
