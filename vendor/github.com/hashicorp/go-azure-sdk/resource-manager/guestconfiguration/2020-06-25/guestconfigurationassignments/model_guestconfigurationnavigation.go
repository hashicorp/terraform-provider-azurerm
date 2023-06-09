package guestconfigurationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GuestConfigurationNavigation struct {
	AssignmentType                  *AssignmentType           `json:"assignmentType,omitempty"`
	ConfigurationParameter          *[]ConfigurationParameter `json:"configurationParameter,omitempty"`
	ConfigurationProtectedParameter *[]ConfigurationParameter `json:"configurationProtectedParameter,omitempty"`
	ConfigurationSetting            *ConfigurationSetting     `json:"configurationSetting,omitempty"`
	ContentHash                     *string                   `json:"contentHash,omitempty"`
	ContentType                     *string                   `json:"contentType,omitempty"`
	ContentUri                      *string                   `json:"contentUri,omitempty"`
	Kind                            *Kind                     `json:"kind,omitempty"`
	Name                            *string                   `json:"name,omitempty"`
	Version                         *string                   `json:"version,omitempty"`
}
