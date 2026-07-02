package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkTrace struct {
	Message *string `json:"message,omitempty"`
	Path    *string `json:"path,omitempty"`
	Status  *string `json:"status,omitempty"`
}
