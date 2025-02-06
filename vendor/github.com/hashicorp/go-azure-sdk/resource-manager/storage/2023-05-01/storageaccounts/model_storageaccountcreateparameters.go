package storageaccounts

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountCreateParameters struct {
	ExtendedLocation *edgezones.Model                          `json:"extendedLocation,omitempty"`
	Identity         *identity.LegacySystemAndUserAssignedMap  `json:"identity,omitempty"`
	Kind             Kind                                      `json:"kind"`
	Location         string                                    `json:"location"`
	Properties       *StorageAccountPropertiesCreateParameters `json:"properties,omitempty"`
	Sku              Sku                                       `json:"sku"`
	Tags             *map[string]string                        `json:"tags,omitempty"`
}
