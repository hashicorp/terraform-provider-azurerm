package restorepoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WinRMConfiguration struct {
	Listeners *[]WinRMListener `json:"listeners,omitempty"`
}
