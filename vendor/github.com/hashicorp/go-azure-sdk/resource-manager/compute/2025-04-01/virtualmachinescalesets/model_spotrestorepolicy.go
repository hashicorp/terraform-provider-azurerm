package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpotRestorePolicy struct {
	Enabled        *bool   `json:"enabled,omitempty"`
	RestoreTimeout *string `json:"restoreTimeout,omitempty"`
}
