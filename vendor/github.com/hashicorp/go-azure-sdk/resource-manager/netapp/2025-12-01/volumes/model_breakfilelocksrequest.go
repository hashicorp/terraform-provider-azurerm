package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BreakFileLocksRequest struct {
	ClientIP                          *string `json:"clientIp,omitempty"`
	ConfirmRunningDisruptiveOperation *bool   `json:"confirmRunningDisruptiveOperation,omitempty"`
}
