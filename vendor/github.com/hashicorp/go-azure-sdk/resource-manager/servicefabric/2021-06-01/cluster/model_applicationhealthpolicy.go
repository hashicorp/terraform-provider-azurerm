package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationHealthPolicy struct {
	DefaultServiceTypeHealthPolicy *ServiceTypeHealthPolicy            `json:"defaultServiceTypeHealthPolicy,omitempty"`
	ServiceTypeHealthPolicies      *map[string]ServiceTypeHealthPolicy `json:"serviceTypeHealthPolicies,omitempty"`
}
