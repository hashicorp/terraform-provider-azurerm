package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableProvidersListCity struct {
	CityName  *string   `json:"cityName,omitempty"`
	Providers *[]string `json:"providers,omitempty"`
}
