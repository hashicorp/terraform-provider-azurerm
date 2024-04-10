package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStorageAccount struct {
	AccountKey         *string             `json:"accountKey,omitempty"`
	AccountName        *string             `json:"accountName,omitempty"`
	AuthenticationMode *AuthenticationMode `json:"authenticationMode,omitempty"`
}
