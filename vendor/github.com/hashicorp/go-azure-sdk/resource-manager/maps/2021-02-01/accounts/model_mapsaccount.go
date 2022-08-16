package accounts

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MapsAccount struct {
	Id         *string                `json:"id,omitempty"`
	Kind       *Kind                  `json:"kind,omitempty"`
	Location   string                 `json:"location"`
	Name       *string                `json:"name,omitempty"`
	Properties *MapsAccountProperties `json:"properties,omitempty"`
	Sku        Sku                    `json:"sku"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Tags       *map[string]string     `json:"tags,omitempty"`
	Type       *string                `json:"type,omitempty"`
}
