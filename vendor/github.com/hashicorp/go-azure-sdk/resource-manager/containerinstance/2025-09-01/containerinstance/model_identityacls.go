package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityAcls struct {
	Acls          *[]IdentityAccessControl `json:"acls,omitempty"`
	DefaultAccess *IdentityAccessLevel     `json:"defaultAccess,omitempty"`
}
