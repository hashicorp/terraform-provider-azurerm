package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Template struct {
	Containers                    *[]Container     `json:"containers,omitempty"`
	InitContainers                *[]BaseContainer `json:"initContainers,omitempty"`
	RevisionSuffix                *string          `json:"revisionSuffix,omitempty"`
	Scale                         *Scale           `json:"scale,omitempty"`
	ServiceBinds                  *[]ServiceBind   `json:"serviceBinds,omitempty"`
	TerminationGracePeriodSeconds *int64           `json:"terminationGracePeriodSeconds,omitempty"`
	Volumes                       *[]Volume        `json:"volumes,omitempty"`
}
