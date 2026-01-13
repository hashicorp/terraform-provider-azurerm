package managedgrafanas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Users struct {
	EditorsCanAdmin *bool `json:"editorsCanAdmin,omitempty"`
	ViewersCanEdit  *bool `json:"viewersCanEdit,omitempty"`
}
