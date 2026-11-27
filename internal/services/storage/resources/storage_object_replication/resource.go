// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_object_replication

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// TODO: @tombuildsstuff: this wants a state migration to move the ID to `{id1}|{id2}` to match other resources

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageObjectReplicationCreate,
		Read:   resourceStorageObjectReplicationRead,
		Update: resourceStorageObjectReplicationUpdate,
		Delete: resourceStorageObjectReplicationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ObjectReplicationID(id)
			return err
		}),

		Schema: storageObjectReplicationSchema(),
	}
}
