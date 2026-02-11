package credential

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CredentialCreateOrUpdateProperties struct {
	Description *string `json:"description,omitempty"`
	Password    string  `json:"password"`
	UserName    string  `json:"userName"`
}
