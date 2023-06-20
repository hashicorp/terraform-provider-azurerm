package apps

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppPatch struct {
	Identity   *identity.SystemAssigned `json:"identity,omitempty"`
	Properties *AppProperties           `json:"properties,omitempty"`
	Sku        *AppSkuInfo              `json:"sku,omitempty"`
	Tags       *map[string]string       `json:"tags,omitempty"`
}
