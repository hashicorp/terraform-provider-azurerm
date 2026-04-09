package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationDeltaHealthPolicy struct {
	DefaultServiceTypeDeltaHealthPolicy *ServiceTypeDeltaHealthPolicy            `json:"defaultServiceTypeDeltaHealthPolicy,omitempty"`
	ServiceTypeDeltaHealthPolicies      *map[string]ServiceTypeDeltaHealthPolicy `json:"serviceTypeDeltaHealthPolicies,omitempty"`
}
