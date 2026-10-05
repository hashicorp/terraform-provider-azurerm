package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemoteDebugging struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Port    *int64 `json:"port,omitempty"`
}
