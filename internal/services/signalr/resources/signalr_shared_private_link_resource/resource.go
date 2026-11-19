// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_shared_private_link_resource

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSignalRSharedPrivateLinkCreateUpdate,
		Read:   resourceSignalRSharedPrivateLinkRead,
		Update: resourceSignalRSharedPrivateLinkCreateUpdate,
		Delete: resourceSignalrSharedPrivateLinkDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := signalr.ParseSharedPrivateLinkResourceIDInsensitively(id)
			return err
		}),

		Schema: signalrSharedPrivateLinkResourceSchema(),
	}
}
