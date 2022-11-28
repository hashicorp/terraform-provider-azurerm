package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MapsAccountUpdateParameters struct {
	Kind       *Kind                  `json:"kind,omitempty"`
	Properties *MapsAccountProperties `json:"properties"`
	Sku        *Sku                   `json:"sku"`
	Tags       *map[string]string     `json:"tags,omitempty"`
}
