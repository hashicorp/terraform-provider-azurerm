package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureKeyVaultLinkedServiceTypeProperties struct {
	BaseURL    interface{}          `json:"baseUrl"`
	Credential *CredentialReference `json:"credential,omitempty"`
}
