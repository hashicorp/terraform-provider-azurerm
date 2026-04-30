package managedinstanceencryptionprotectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstanceEncryptionProtectorProperties struct {
	AutoRotationEnabled *bool         `json:"autoRotationEnabled,omitempty"`
	ServerKeyName       *string       `json:"serverKeyName,omitempty"`
	ServerKeyType       ServerKeyType `json:"serverKeyType"`
	Thumbprint          *string       `json:"thumbprint,omitempty"`
	Uri                 *string       `json:"uri,omitempty"`
}
