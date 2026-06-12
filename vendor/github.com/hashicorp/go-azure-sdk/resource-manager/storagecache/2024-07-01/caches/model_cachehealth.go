package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheHealth struct {
	Conditions        *[]Condition     `json:"conditions,omitempty"`
	State             *HealthStateType `json:"state,omitempty"`
	StatusDescription *string          `json:"statusDescription,omitempty"`
}
