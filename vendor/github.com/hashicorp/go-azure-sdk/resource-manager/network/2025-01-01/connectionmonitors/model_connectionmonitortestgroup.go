package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMonitorTestGroup struct {
	Destinations       []string `json:"destinations"`
	Disable            *bool    `json:"disable,omitempty"`
	Name               string   `json:"name"`
	Sources            []string `json:"sources"`
	TestConfigurations []string `json:"testConfigurations"`
}
