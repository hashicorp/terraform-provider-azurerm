package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileShare struct {
	Password *string `json:"password,omitempty"`
	Path     string  `json:"path"`
	UserName *string `json:"userName,omitempty"`
}
