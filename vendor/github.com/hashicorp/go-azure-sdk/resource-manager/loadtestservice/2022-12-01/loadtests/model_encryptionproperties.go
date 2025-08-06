package loadtests

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionProperties struct {
	Identity *EncryptionPropertiesIdentity `json:"identity,omitempty"`
	KeyURL   *string                       `json:"keyUrl,omitempty"`
}
