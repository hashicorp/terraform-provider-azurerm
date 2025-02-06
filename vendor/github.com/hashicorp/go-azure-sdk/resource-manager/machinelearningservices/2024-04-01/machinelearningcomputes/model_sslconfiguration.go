package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SslConfiguration struct {
	Cert                    *string          `json:"cert,omitempty"`
	Cname                   *string          `json:"cname,omitempty"`
	Key                     *string          `json:"key,omitempty"`
	LeafDomainLabel         *string          `json:"leafDomainLabel,omitempty"`
	OverwriteExistingDomain *bool            `json:"overwriteExistingDomain,omitempty"`
	Status                  *SslConfigStatus `json:"status,omitempty"`
}
