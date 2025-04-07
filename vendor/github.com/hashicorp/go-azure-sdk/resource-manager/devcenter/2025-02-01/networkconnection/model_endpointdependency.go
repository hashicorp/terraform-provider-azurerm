package networkconnection

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointDependency struct {
	Description     *string           `json:"description,omitempty"`
	DomainName      *string           `json:"domainName,omitempty"`
	EndpointDetails *[]EndpointDetail `json:"endpointDetails,omitempty"`
}
