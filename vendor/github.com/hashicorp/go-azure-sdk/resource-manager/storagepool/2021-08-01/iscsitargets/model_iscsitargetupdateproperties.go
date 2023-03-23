package iscsitargets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IscsiTargetUpdateProperties struct {
	Luns       *[]IscsiLun `json:"luns,omitempty"`
	StaticAcls *[]Acl      `json:"staticAcls,omitempty"`
}
