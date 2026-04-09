package deviceupdates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointUpdate struct {
	Id                      *string `json:"id,omitempty"`
	ImmutableResourceId     *string `json:"immutableResourceId,omitempty"`
	ImmutableSubscriptionId *string `json:"immutableSubscriptionId,omitempty"`
	Location                *string `json:"location,omitempty"`
	VnetTrafficTag          *string `json:"vnetTrafficTag,omitempty"`
}
