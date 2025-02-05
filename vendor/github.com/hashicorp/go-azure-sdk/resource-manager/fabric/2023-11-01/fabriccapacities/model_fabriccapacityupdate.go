package fabriccapacities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricCapacityUpdate struct {
	Properties *FabricCapacityUpdateProperties `json:"properties,omitempty"`
	Sku        *RpSku                          `json:"sku,omitempty"`
	Tags       *map[string]string              `json:"tags,omitempty"`
}
