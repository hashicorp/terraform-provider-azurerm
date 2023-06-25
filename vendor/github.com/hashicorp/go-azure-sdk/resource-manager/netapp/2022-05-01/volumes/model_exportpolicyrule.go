package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportPolicyRule struct {
	AllowedClients      *string    `json:"allowedClients,omitempty"`
	ChownMode           *ChownMode `json:"chownMode,omitempty"`
	Cifs                *bool      `json:"cifs,omitempty"`
	HasRootAccess       *bool      `json:"hasRootAccess,omitempty"`
	Kerberos5ReadOnly   *bool      `json:"kerberos5ReadOnly,omitempty"`
	Kerberos5ReadWrite  *bool      `json:"kerberos5ReadWrite,omitempty"`
	Kerberos5iReadOnly  *bool      `json:"kerberos5iReadOnly,omitempty"`
	Kerberos5iReadWrite *bool      `json:"kerberos5iReadWrite,omitempty"`
	Kerberos5pReadOnly  *bool      `json:"kerberos5pReadOnly,omitempty"`
	Kerberos5pReadWrite *bool      `json:"kerberos5pReadWrite,omitempty"`
	Nfsv3               *bool      `json:"nfsv3,omitempty"`
	Nfsv41              *bool      `json:"nfsv41,omitempty"`
	RuleIndex           *int64     `json:"ruleIndex,omitempty"`
	UnixReadOnly        *bool      `json:"unixReadOnly,omitempty"`
	UnixReadWrite       *bool      `json:"unixReadWrite,omitempty"`
}
