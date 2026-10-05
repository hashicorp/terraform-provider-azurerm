package account

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountSku struct {
	Capacity *int64 `json:"capacity,omitempty"`
	Name     *Name  `json:"name,omitempty"`
}
