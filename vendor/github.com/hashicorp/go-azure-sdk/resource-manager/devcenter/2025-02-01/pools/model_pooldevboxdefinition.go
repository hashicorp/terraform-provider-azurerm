package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolDevBoxDefinition struct {
	ActiveImageReference *ImageReference `json:"activeImageReference,omitempty"`
	ImageReference       *ImageReference `json:"imageReference,omitempty"`
	Sku                  *Sku            `json:"sku,omitempty"`
}
