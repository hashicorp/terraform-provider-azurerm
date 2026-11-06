package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSImageNotificationProfile struct {
	Enable           *bool   `json:"enable,omitempty"`
	NotBeforeTimeout *string `json:"notBeforeTimeout,omitempty"`
}
