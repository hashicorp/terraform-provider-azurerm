package signalr

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveTraceConfiguration struct {
	Categories *[]LiveTraceCategory `json:"categories,omitempty"`
	Enabled    *string              `json:"enabled,omitempty"`
}
