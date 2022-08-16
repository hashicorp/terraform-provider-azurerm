package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthMonitor struct {
	Details *string `json:"details,omitempty"`
	Id      *string `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
}
