package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceNameAvailabilityRequest struct {
	EnvironmentId *string                `json:"environmentId,omitempty"`
	IsFqdn        *bool                  `json:"isFqdn,omitempty"`
	Name          string                 `json:"name"`
	Type          CheckNameResourceTypes `json:"type"`
}
