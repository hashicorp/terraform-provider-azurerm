package alertrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IncidentConfiguration struct {
	CreateIncident        bool                   `json:"createIncident"`
	GroupingConfiguration *GroupingConfiguration `json:"groupingConfiguration,omitempty"`
}
