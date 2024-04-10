package staticsites

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteARMResource struct {
	Id         *string                            `json:"id,omitempty"`
	Identity   *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Kind       *string                            `json:"kind,omitempty"`
	Location   string                             `json:"location"`
	Name       *string                            `json:"name,omitempty"`
	Properties *StaticSite                        `json:"properties,omitempty"`
	Sku        *SkuDescription                    `json:"sku,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
	Type       *string                            `json:"type,omitempty"`
}
