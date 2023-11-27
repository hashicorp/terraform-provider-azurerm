package databases

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Database struct {
	Id         *string                   `json:"id,omitempty"`
	Identity   *identity.UserAssignedMap `json:"identity,omitempty"`
	Kind       *string                   `json:"kind,omitempty"`
	Location   string                    `json:"location"`
	ManagedBy  *string                   `json:"managedBy,omitempty"`
	Name       *string                   `json:"name,omitempty"`
	Properties *DatabaseProperties       `json:"properties,omitempty"`
	Sku        *Sku                      `json:"sku,omitempty"`
	Tags       *map[string]string        `json:"tags,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}
