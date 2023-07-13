package managedhsms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHsmSku struct {
	Family ManagedHsmSkuFamily `json:"family"`
	Name   ManagedHsmSkuName   `json:"name"`
}
