// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ssh_public_key

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/sshpublickeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSshPublicKeyCreate,
		Read:   resourceSshPublicKeyRead,
		Update: resourceSshPublicKeyUpdate,
		Delete: resourceSshPublicKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := sshpublickeys.ParseSshPublicKeyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(45 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
		},

		Schema: sshPublicKeySchema(),
	}
}
