package volumesreplication

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BreakReplicationRequest struct {
	ForceBreakReplication *bool `json:"forceBreakReplication,omitempty"`
}
