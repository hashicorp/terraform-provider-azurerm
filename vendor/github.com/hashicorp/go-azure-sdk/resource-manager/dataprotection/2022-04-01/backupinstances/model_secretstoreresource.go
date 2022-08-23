package backupinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretStoreResource struct {
	SecretStoreType SecretStoreType `json:"secretStoreType"`
	Uri             *string         `json:"uri,omitempty"`
	Value           *string         `json:"value,omitempty"`
}
