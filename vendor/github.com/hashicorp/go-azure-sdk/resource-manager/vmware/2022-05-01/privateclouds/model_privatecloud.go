package privateclouds

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateCloud struct {
	Id         *string                  `json:"id,omitempty"`
	Identity   *identity.SystemAssigned `json:"identity,omitempty"`
	Location   *string                  `json:"location,omitempty"`
	Name       *string                  `json:"name,omitempty"`
	Properties *PrivateCloudProperties  `json:"properties,omitempty"`
	Sku        Sku                      `json:"sku"`
	Tags       *map[string]string       `json:"tags,omitempty"`
	Type       *string                  `json:"type,omitempty"`
}
