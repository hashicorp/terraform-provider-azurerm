package subnets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceEndpointPolicy struct {
	Etag       *string                                `json:"etag,omitempty"`
	Id         *string                                `json:"id,omitempty"`
	Kind       *string                                `json:"kind,omitempty"`
	Location   *string                                `json:"location,omitempty"`
	Name       *string                                `json:"name,omitempty"`
	Properties *ServiceEndpointPolicyPropertiesFormat `json:"properties,omitempty"`
	Tags       *map[string]string                     `json:"tags,omitempty"`
	Type       *string                                `json:"type,omitempty"`
}
