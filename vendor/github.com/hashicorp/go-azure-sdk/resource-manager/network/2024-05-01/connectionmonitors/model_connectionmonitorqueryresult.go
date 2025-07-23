package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMonitorQueryResult struct {
	SourceStatus *ConnectionMonitorSourceStatus `json:"sourceStatus,omitempty"`
	States       *[]ConnectionStateSnapshot     `json:"states,omitempty"`
}
