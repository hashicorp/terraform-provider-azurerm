package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespacePatchParameters struct {
	Sku  *Sku               `json:"sku,omitempty"`
	Tags *map[string]string `json:"tags,omitempty"`
}
