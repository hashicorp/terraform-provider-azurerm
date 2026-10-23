package autonomousdatabases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionStringType struct {
	AllConnectionStrings *AllConnectionStringType `json:"allConnectionStrings,omitempty"`
	Dedicated            *string                  `json:"dedicated,omitempty"`
	High                 *string                  `json:"high,omitempty"`
	Low                  *string                  `json:"low,omitempty"`
	Medium               *string                  `json:"medium,omitempty"`
	Profiles             *[]ProfileType           `json:"profiles,omitempty"`
}
