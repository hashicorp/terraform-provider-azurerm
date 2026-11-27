// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_network_acl

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages the Network ACL for a Web Pubsub.",
		Create:      resourceWebPubsubNetworkACLCreateUpdate,
		Read:        resourceWebPubsubNetworkACLRead,
		Update:      resourceWebPubsubNetworkACLCreateUpdate,
		Delete:      resourceWebpubsubNetworkACLDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webpubsub.ParseWebPubSubID(id)
			return err
		}),

		Schema: webPubsubNetworkAclSchema(),
	}
}
