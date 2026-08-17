package healthbots

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthBotUpdateParameters struct {
	Identity   *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Location   *string                            `json:"location,omitempty"`
	Properties *HealthBotProperties               `json:"properties,omitempty"`
	Sku        *Sku                               `json:"sku,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
}
