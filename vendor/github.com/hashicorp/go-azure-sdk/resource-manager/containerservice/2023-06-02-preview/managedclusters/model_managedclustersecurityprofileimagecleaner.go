package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterSecurityProfileImageCleaner struct {
	Enabled       *bool  `json:"enabled,omitempty"`
	IntervalHours *int64 `json:"intervalHours,omitempty"`
}
