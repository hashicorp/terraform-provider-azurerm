package serverrestart

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerRestartParameter struct {
	MaxFailoverSeconds  *int64            `json:"maxFailoverSeconds,omitempty"`
	RestartWithFailover *EnableStatusEnum `json:"restartWithFailover,omitempty"`
}
