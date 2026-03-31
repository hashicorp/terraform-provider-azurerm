package tenants

// Copyright IBM Corp. 2021, 2025 All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Sku struct {
	Name SkuName `json:"name"`
	Tier SkuTier `json:"tier"`
}
