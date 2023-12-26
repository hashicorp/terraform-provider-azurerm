package saplandscapemonitor

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapLandscapeMonitorListResult struct {
	NextLink *string                `json:"nextLink,omitempty"`
	Value    *[]SapLandscapeMonitor `json:"value,omitempty"`
}
