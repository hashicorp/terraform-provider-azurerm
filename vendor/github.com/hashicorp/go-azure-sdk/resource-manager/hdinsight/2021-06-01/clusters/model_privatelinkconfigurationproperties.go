package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkConfigurationProperties struct {
	GroupId           string                                     `json:"groupId"`
	IPConfigurations  []IPConfiguration                          `json:"ipConfigurations"`
	ProvisioningState *PrivateLinkConfigurationProvisioningState `json:"provisioningState,omitempty"`
}
