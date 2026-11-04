package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Sku struct {
	Capacity *int64  `json:"capacity,omitempty"`
	Name     *string `json:"name,omitempty"`
	Tier     *string `json:"tier,omitempty"`
}
