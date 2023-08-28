package assignment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Assignment struct {
	Id         *string                          `json:"id,omitempty"`
	Identity   identity.SystemOrUserAssignedMap `json:"identity"`
	Location   string                           `json:"location"`
	Name       *string                          `json:"name,omitempty"`
	Properties AssignmentProperties             `json:"properties"`
	Type       *string                          `json:"type,omitempty"`
}
