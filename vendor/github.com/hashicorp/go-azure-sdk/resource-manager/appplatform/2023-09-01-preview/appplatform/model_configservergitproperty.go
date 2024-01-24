package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigServerGitProperty struct {
	HostKey               *string                 `json:"hostKey,omitempty"`
	HostKeyAlgorithm      *string                 `json:"hostKeyAlgorithm,omitempty"`
	Label                 *string                 `json:"label,omitempty"`
	Password              *string                 `json:"password,omitempty"`
	PrivateKey            *string                 `json:"privateKey,omitempty"`
	Repositories          *[]GitPatternRepository `json:"repositories,omitempty"`
	SearchPaths           *[]string               `json:"searchPaths,omitempty"`
	StrictHostKeyChecking *bool                   `json:"strictHostKeyChecking,omitempty"`
	Uri                   string                  `json:"uri"`
	Username              *string                 `json:"username,omitempty"`
}
