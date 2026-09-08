package mongorbacs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivilegeResource struct {
	Collection *string `json:"collection,omitempty"`
	Db         *string `json:"db,omitempty"`
}
