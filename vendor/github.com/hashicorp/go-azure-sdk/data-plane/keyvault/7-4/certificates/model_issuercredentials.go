package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IssuerCredentials struct {
	AccountId *string `json:"account_id,omitempty"`
	Pwd       *string `json:"pwd,omitempty"`
}
