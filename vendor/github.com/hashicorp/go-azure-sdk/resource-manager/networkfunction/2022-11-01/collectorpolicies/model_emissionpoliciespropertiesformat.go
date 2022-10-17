package collectorpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EmissionPoliciesPropertiesFormat struct {
	EmissionDestinations *[]EmissionPolicyDestination `json:"emissionDestinations,omitempty"`
	EmissionType         *EmissionType                `json:"emissionType,omitempty"`
}
