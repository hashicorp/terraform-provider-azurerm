package backend

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendProxyContract struct {
	Password *string `json:"password,omitempty"`
	Url      string  `json:"url"`
	Username *string `json:"username,omitempty"`
}
