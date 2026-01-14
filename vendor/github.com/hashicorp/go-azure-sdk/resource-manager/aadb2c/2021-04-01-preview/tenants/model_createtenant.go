package tenants

// Copyright IBM Corp. 2021, 2025 All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateTenant struct {
	Location   Location                  `json:"location"`
	Properties TenantPropertiesForCreate `json:"properties"`
	Sku        Sku                       `json:"sku"`
	Tags       *map[string]string        `json:"tags,omitempty"`
}
