package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourcesResponseEndpointsItem struct {
	CustomDomains *[]ResourcesResponseEndpointsPropertiesItemsItem `json:"customDomains,omitempty"`
	History       *bool                                            `json:"history,omitempty"`
	Id            *string                                          `json:"id,omitempty"`
	Name          *string                                          `json:"name,omitempty"`
}
