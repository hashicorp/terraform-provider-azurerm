package serverrestart

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestartParameter struct {
	FailoverMode        *FailoverMode `json:"failoverMode,omitempty"`
	RestartWithFailover *bool         `json:"restartWithFailover,omitempty"`
}
