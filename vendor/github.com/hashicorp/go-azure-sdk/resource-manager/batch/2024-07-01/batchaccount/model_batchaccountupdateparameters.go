package batchaccount

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchAccountUpdateParameters struct {
	Identity   *identity.SystemOrUserAssignedMap `json:"identity,omitempty"`
	Properties *BatchAccountUpdateProperties     `json:"properties,omitempty"`
	Tags       *map[string]string                `json:"tags,omitempty"`
}
