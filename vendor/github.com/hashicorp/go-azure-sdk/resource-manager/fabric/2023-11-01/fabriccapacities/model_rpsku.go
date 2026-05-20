package fabriccapacities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RpSku struct {
	Name string    `json:"name"`
	Tier RpSkuTier `json:"tier"`
}
