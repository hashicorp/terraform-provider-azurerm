package fluidrelayservers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomerManagedKeyEncryptionProperties struct {
	KeyEncryptionKeyIdentity *CustomerManagedKeyEncryptionPropertiesKeyEncryptionKeyIdentity `json:"keyEncryptionKeyIdentity,omitempty"`
	KeyEncryptionKeyUrl      *string                                                         `json:"keyEncryptionKeyUrl,omitempty"`
}
