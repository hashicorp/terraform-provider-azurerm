package groupuser

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserIdentityContract struct {
	Id       *string `json:"id,omitempty"`
	Provider *string `json:"provider,omitempty"`
}
