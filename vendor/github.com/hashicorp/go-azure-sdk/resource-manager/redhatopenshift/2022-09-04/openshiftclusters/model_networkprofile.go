package openshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProfile struct {
	PodCidr     *string `json:"podCidr,omitempty"`
	ServiceCidr *string `json:"serviceCidr,omitempty"`
}
