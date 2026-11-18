package networkprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointIPConfiguration struct {
	Etag       *string                                   `json:"etag,omitempty"`
	Name       *string                                   `json:"name,omitempty"`
	Properties *PrivateEndpointIPConfigurationProperties `json:"properties,omitempty"`
	Type       *string                                   `json:"type,omitempty"`
}
