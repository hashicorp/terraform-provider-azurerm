// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageIDValidationFunc func(id, storageDomainSuffix string) error

// ImporterValidatingStorageResourceId is similar to pluginsdk.ImporterValidatingResourceId whilst additionally passing in
// the storage domain suffix for the current environment/cloud, so that the domain component of the ID can be validated.
func ImporterValidatingStorageResourceId(validateFunc StorageIDValidationFunc) *schema.ResourceImporter {
	return ImporterValidatingStorageResourceIdThen(validateFunc, nil)
}

// ImporterValidatingStorageResourceIdThen is similar to pluginsdk.ImporterValidatingResourceIdThen whilst additionally passing in
// the storage domain suffix for the current environment/cloud, so that the domain component of the ID can be validated.
func ImporterValidatingStorageResourceIdThen(validateFunc StorageIDValidationFunc, thenFunc pluginsdk.ImporterFunc) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			storageDomainSuffix := meta.(*clients.Client).Storage.StorageDomainSuffix
			log.Printf("[DEBUG] Importing Storage Resource - parsing %q using domain suffix %q", d.Id(), storageDomainSuffix)
			if err := validateFunc(d.Id(), storageDomainSuffix); err != nil {
				return []*pluginsdk.ResourceData{d}, err
			}

			if thenFunc != nil {
				return thenFunc(ctx, d, meta)
			}

			return []*pluginsdk.ResourceData{d}, nil
		},
	}
}
