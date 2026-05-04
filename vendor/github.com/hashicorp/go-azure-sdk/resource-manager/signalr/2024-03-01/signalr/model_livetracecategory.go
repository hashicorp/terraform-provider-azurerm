package signalr

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveTraceCategory struct {
	Enabled *string `json:"enabled,omitempty"`
	Name    *string `json:"name,omitempty"`
}
