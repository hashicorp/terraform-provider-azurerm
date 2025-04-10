package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WhatIfChange struct {
	After             *interface{}            `json:"after,omitempty"`
	Before            *interface{}            `json:"before,omitempty"`
	ChangeType        ChangeType              `json:"changeType"`
	Delta             *[]WhatIfPropertyChange `json:"delta,omitempty"`
	ResourceId        string                  `json:"resourceId"`
	UnsupportedReason *string                 `json:"unsupportedReason,omitempty"`
}
