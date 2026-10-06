package backend

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendPoolItem struct {
	Id       string `json:"id"`
	Priority *int64 `json:"priority,omitempty"`
	Weight   *int64 `json:"weight,omitempty"`
}
