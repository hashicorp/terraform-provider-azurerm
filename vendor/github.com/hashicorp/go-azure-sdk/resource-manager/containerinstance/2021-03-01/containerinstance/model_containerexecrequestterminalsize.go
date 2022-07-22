package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerExecRequestTerminalSize struct {
	Cols *int64 `json:"cols,omitempty"`
	Rows *int64 `json:"rows,omitempty"`
}
