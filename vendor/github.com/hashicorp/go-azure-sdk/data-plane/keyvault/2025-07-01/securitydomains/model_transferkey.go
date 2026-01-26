package securitydomains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransferKey struct {
	KeyFormat   *string                  `json:"key_format,omitempty"`
	TransferKey SecurityDomainJsonWebKey `json:"transfer_key"`
}
