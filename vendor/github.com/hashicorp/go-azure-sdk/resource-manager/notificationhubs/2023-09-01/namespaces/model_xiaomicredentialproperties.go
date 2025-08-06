package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type XiaomiCredentialProperties struct {
	AppSecret *string `json:"appSecret,omitempty"`
	Endpoint  *string `json:"endpoint,omitempty"`
}
