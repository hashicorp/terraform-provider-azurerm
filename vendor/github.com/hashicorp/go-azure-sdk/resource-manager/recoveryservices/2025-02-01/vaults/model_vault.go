package vaults

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Vault struct {
	Etag       *string                            `json:"etag,omitempty"`
	Id         *string                            `json:"id,omitempty"`
	Identity   *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Location   string                             `json:"location"`
	Name       *string                            `json:"name,omitempty"`
	Properties *VaultProperties                   `json:"properties,omitempty"`
	Sku        *Sku                               `json:"sku,omitempty"`
	SystemData *systemdata.SystemData             `json:"systemData,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
	Type       *string                            `json:"type,omitempty"`
}
