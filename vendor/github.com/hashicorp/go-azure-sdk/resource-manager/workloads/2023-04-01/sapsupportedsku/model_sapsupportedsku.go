package sapsupportedsku

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPSupportedSku struct {
	IsAppServerCertified *bool   `json:"isAppServerCertified,omitempty"`
	IsDatabaseCertified  *bool   `json:"isDatabaseCertified,omitempty"`
	VMSku                *string `json:"vmSku,omitempty"`
}
