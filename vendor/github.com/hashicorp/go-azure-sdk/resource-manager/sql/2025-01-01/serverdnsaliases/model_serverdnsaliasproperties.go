package serverdnsaliases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerDnsAliasProperties struct {
	AzureDnsRecord *string `json:"azureDnsRecord,omitempty"`
}
