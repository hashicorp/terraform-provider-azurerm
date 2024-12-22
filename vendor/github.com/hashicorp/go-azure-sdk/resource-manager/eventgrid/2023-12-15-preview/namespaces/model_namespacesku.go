package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespaceSku struct {
	Capacity *int64   `json:"capacity,omitempty"`
	Name     *SkuName `json:"name,omitempty"`
}
