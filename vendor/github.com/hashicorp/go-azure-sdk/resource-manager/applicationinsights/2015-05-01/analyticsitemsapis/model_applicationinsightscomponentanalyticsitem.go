package analyticsitemsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentAnalyticsItem struct {
	Content      *string                                              `json:"Content,omitempty"`
	Id           *string                                              `json:"Id,omitempty"`
	Name         *string                                              `json:"Name,omitempty"`
	Properties   *ApplicationInsightsComponentAnalyticsItemProperties `json:"Properties,omitempty"`
	Scope        *ItemScope                                           `json:"Scope,omitempty"`
	TimeCreated  *string                                              `json:"TimeCreated,omitempty"`
	TimeModified *string                                              `json:"TimeModified,omitempty"`
	Type         *ItemType                                            `json:"Type,omitempty"`
	Version      *string                                              `json:"Version,omitempty"`
}
