package tenants

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Sku struct {
	Name SkuName `json:"name"`
	Tier SkuTier `json:"tier"`
}
