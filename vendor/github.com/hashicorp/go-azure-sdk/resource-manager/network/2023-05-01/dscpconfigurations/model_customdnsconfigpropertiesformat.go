package dscpconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDnsConfigPropertiesFormat struct {
	Fqdn        *string   `json:"fqdn,omitempty"`
	IPAddresses *[]string `json:"ipAddresses,omitempty"`
}
