package pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Pool struct {
	Etag       *string                   `json:"etag,omitempty"`
	Id         *string                   `json:"id,omitempty"`
	Identity   *identity.UserAssignedMap `json:"identity,omitempty"`
	Name       *string                   `json:"name,omitempty"`
	Properties *PoolProperties           `json:"properties,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}
