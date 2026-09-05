package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticTrafficFilter struct {
	Description      *string                     `json:"description,omitempty"`
	Id               *string                     `json:"id,omitempty"`
	IncludeByDefault *bool                       `json:"includeByDefault,omitempty"`
	Name             *string                     `json:"name,omitempty"`
	Region           *string                     `json:"region,omitempty"`
	Rules            *[]ElasticTrafficFilterRule `json:"rules,omitempty"`
	Type             *Type                       `json:"type,omitempty"`
}
