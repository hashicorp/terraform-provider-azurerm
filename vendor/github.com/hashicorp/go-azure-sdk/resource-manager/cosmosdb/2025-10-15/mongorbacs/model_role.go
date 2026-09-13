package mongorbacs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Role struct {
	Db   *string `json:"db,omitempty"`
	Role *string `json:"role,omitempty"`
}
