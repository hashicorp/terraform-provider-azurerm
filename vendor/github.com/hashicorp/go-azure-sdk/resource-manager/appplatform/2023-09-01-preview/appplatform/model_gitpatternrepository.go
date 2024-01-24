package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GitPatternRepository struct {
	HostKey               *string   `json:"hostKey,omitempty"`
	HostKeyAlgorithm      *string   `json:"hostKeyAlgorithm,omitempty"`
	Label                 *string   `json:"label,omitempty"`
	Name                  string    `json:"name"`
	Password              *string   `json:"password,omitempty"`
	Pattern               *[]string `json:"pattern,omitempty"`
	PrivateKey            *string   `json:"privateKey,omitempty"`
	SearchPaths           *[]string `json:"searchPaths,omitempty"`
	StrictHostKeyChecking *bool     `json:"strictHostKeyChecking,omitempty"`
	Uri                   string    `json:"uri"`
	Username              *string   `json:"username,omitempty"`
}
