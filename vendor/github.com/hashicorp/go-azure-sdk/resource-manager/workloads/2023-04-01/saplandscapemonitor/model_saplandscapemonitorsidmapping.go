package saplandscapemonitor

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapLandscapeMonitorSidMapping struct {
	Name   *string   `json:"name,omitempty"`
	TopSid *[]string `json:"topSid,omitempty"`
}
