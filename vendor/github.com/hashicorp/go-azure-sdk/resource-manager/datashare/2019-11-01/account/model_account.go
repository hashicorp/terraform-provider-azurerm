package account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Account struct {
	Id         *string                 `json:"id,omitempty"`
	Identity   identity.SystemAssigned `json:"identity"`
	Location   *string                 `json:"location,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *AccountProperties      `json:"properties,omitempty"`
	Tags       *map[string]string      `json:"tags,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}
