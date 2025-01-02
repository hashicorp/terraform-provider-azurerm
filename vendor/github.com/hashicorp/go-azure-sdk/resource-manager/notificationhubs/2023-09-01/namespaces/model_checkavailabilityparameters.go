package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckAvailabilityParameters struct {
	Id           *string            `json:"id,omitempty"`
	IsAvailiable *bool              `json:"isAvailiable,omitempty"`
	Location     *string            `json:"location,omitempty"`
	Name         string             `json:"name"`
	Sku          *Sku               `json:"sku,omitempty"`
	Tags         *map[string]string `json:"tags,omitempty"`
	Type         *string            `json:"type,omitempty"`
}
