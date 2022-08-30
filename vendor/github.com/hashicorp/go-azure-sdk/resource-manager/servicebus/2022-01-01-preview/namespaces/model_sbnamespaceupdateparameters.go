package namespaces

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SBNamespaceUpdateParameters struct {
	Id         *string                            `json:"id,omitempty"`
	Identity   *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Location   *string                            `json:"location,omitempty"`
	Name       *string                            `json:"name,omitempty"`
	Properties *SBNamespaceUpdateProperties       `json:"properties,omitempty"`
	Sku        *SBSku                             `json:"sku,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
	Type       *string                            `json:"type,omitempty"`
}
