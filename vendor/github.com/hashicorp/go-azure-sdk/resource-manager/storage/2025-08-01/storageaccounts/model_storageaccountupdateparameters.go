package storageaccounts

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountUpdateParameters struct {
	Identity   *identity.LegacySystemAndUserAssignedMap  `json:"identity,omitempty"`
	Kind       *Kind                                     `json:"kind,omitempty"`
	Placement  *Placement                                `json:"placement,omitempty"`
	Properties *StorageAccountPropertiesUpdateParameters `json:"properties,omitempty"`
	Sku        *Sku                                      `json:"sku,omitempty"`
	Tags       *map[string]string                        `json:"tags,omitempty"`
	Zones      *zones.Schema                             `json:"zones,omitempty"`
}
