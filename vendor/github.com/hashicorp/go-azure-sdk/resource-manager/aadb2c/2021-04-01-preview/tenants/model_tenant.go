package tenants

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Tenant struct {
	Id         *string            `json:"id,omitempty"`
	Location   *Location          `json:"location,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties *TenantProperties  `json:"properties,omitempty"`
	Sku        *Sku               `json:"sku,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
