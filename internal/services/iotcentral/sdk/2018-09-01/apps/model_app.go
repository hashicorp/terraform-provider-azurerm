package apps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type App struct {
	Id         *string            `json:"id,omitempty"`
	Location   string             `json:"location"`
	Name       *string            `json:"name,omitempty"`
	Properties *AppProperties     `json:"properties,omitempty"`
	Sku        AppSkuInfo         `json:"sku"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
