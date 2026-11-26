package credential

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CredentialUpdateProperties struct {
	Description *string `json:"description,omitempty"`
	Password    *string `json:"password,omitempty"`
	UserName    *string `json:"userName,omitempty"`
}
