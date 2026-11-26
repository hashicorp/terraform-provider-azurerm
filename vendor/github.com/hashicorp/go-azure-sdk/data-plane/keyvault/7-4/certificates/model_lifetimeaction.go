package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LifetimeAction struct {
	Action  *Action  `json:"action,omitempty"`
	Trigger *Trigger `json:"trigger,omitempty"`
}
