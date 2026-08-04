package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LifetimeActionsTrigger struct {
	TimeAfterCreate  *string `json:"timeAfterCreate,omitempty"`
	TimeBeforeExpiry *string `json:"timeBeforeExpiry,omitempty"`
}
