package grafanaresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Smtp struct {
	Enabled        *bool           `json:"enabled,omitempty"`
	FromAddress    *string         `json:"fromAddress,omitempty"`
	FromName       *string         `json:"fromName,omitempty"`
	Host           *string         `json:"host,omitempty"`
	Password       *string         `json:"password,omitempty"`
	SkipVerify     *bool           `json:"skipVerify,omitempty"`
	StartTLSPolicy *StartTLSPolicy `json:"startTLSPolicy,omitempty"`
	User           *string         `json:"user,omitempty"`
}
