package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableProvidersListState struct {
	Cities    *[]AvailableProvidersListCity `json:"cities,omitempty"`
	Providers *[]string                     `json:"providers,omitempty"`
	StateName *string                       `json:"stateName,omitempty"`
}
