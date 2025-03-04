package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrowserCredentialProperties struct {
	Subject         string `json:"subject"`
	VapidPrivateKey string `json:"vapidPrivateKey"`
	VapidPublicKey  string `json:"vapidPublicKey"`
}
