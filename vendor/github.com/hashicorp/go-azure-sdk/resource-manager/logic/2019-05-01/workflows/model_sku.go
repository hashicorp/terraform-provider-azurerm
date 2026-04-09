package workflows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Sku struct {
	Name SkuName            `json:"name"`
	Plan *ResourceReference `json:"plan,omitempty"`
}
