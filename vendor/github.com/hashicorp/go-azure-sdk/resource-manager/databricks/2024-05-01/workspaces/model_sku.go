package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Sku struct {
	Name string  `json:"name"`
	Tier *string `json:"tier,omitempty"`
}
