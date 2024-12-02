package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMonitor struct {
	Location   *string                     `json:"location,omitempty"`
	Properties ConnectionMonitorParameters `json:"properties"`
	Tags       *map[string]string          `json:"tags,omitempty"`
}
