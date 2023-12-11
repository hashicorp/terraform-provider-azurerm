package updates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdatePrerequisite struct {
	PackageName *string `json:"packageName,omitempty"`
	UpdateType  *string `json:"updateType,omitempty"`
	Version     *string `json:"version,omitempty"`
}
