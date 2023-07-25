// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = QueueV0ToV1{}

type QueueV0ToV1 struct{}

func (QueueV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"storage_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (QueueV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	// this should have been applied from pre-0.12 migration system; backporting just in-case
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		environment := meta.(*clients.Client).Account.Environment
		storageDomainSuffix, ok := environment.Storage.DomainSuffix()
		if !ok {
			return nil, fmt.Errorf("could not determine Storage domain suffix for environment %q", environment.Name)
		}

		queueName := rawState["name"]
		storageAccountName := rawState["storage_account_name"]
		newID := fmt.Sprintf("https://%s.queue.%s/%s", storageAccountName, *storageDomainSuffix, queueName)
		rawState["id"] = newID

		return rawState, nil
	}
}
