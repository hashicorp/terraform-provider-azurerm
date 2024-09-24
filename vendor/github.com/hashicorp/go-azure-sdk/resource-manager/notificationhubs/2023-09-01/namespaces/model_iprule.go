package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPRule struct {
	IPMask string         `json:"ipMask"`
	Rights []AccessRights `json:"rights"`
}
