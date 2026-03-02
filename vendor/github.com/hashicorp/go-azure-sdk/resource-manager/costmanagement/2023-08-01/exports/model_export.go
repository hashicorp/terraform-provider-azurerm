package exports

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Export struct {
	ETag       *string                  `json:"eTag,omitempty"`
	Id         *string                  `json:"id,omitempty"`
	Identity   *identity.SystemAssigned `json:"identity,omitempty"`
	Location   *string                  `json:"location,omitempty"`
	Name       *string                  `json:"name,omitempty"`
	Properties *ExportProperties        `json:"properties,omitempty"`
	Type       *string                  `json:"type,omitempty"`
}
