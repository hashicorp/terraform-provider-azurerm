package sourcecontrolconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceControlConfiguration struct {
	Id         *string                               `json:"id,omitempty"`
	Name       *string                               `json:"name,omitempty"`
	Properties *SourceControlConfigurationProperties `json:"properties,omitempty"`
	SystemData *SystemData                           `json:"systemData,omitempty"`
	Type       *string                               `json:"type,omitempty"`
}
