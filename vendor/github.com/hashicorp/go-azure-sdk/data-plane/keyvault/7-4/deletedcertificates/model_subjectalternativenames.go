package deletedcertificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubjectAlternativeNames struct {
	DnsNames *[]string `json:"dns_names,omitempty"`
	Emails   *[]string `json:"emails,omitempty"`
	Upns     *[]string `json:"upns,omitempty"`
}
