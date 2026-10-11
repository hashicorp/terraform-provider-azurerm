package connectedregistries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoginServerProperties struct {
	Host *string        `json:"host,omitempty"`
	Tls  *TlsProperties `json:"tls,omitempty"`
}
