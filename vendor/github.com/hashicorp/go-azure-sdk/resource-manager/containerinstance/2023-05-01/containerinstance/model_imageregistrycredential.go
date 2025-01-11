package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageRegistryCredential struct {
	Identity    *string `json:"identity,omitempty"`
	IdentityURL *string `json:"identityUrl,omitempty"`
	Password    *string `json:"password,omitempty"`
	Server      string  `json:"server"`
	Username    *string `json:"username,omitempty"`
}
