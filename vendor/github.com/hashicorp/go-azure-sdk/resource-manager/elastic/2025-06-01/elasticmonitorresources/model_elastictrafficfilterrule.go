package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticTrafficFilterRule struct {
	AzureEndpointGuid *string `json:"azureEndpointGuid,omitempty"`
	AzureEndpointName *string `json:"azureEndpointName,omitempty"`
	Description       *string `json:"description,omitempty"`
	Id                *string `json:"id,omitempty"`
	Source            *string `json:"source,omitempty"`
}
